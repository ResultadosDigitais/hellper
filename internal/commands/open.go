package commands

import (
	"context"
	"fmt"
	"hellper/internal/app"
	"hellper/internal/concurrence"
	"strconv"
	"strings"
	"sync"
	"time"

	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/meeting"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

// OpenStartIncidentDialog opens a dialog on Slack, so the user can start an incident
func OpenStartIncidentDialog(ctx context.Context, app *app.App, userID string, triggerID string) error {
	services, err := app.ServiceRepository.ListServiceInstances(ctx)
	if err != nil {
		return err
	}

	serviceList := getDialogOptionsWithServiceInstances(services)

	incidentTitle := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident Title",
			Name:        "incident_title",
			Type:        "text",
			Placeholder: "My Incident Title",
		},
		MaxLength: 100,
	}

	product := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Product / Service",
			Name:        "product",
			Type:        "select",
			Placeholder: "Set the product / service",
			Optional:    false,
		},
		Options:      serviceList,
		OptionGroups: []slack.DialogOptionGroup{},
	}

	commander := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Incident commander",
			Name:        "incident_commander",
			Type:        "select",
			Placeholder: "Set the Incident commander",
			Optional:    false,
		},
		Value:        userID,
		DataSource:   "users",
		OptionGroups: []slack.DialogOptionGroup{},
	}

	severityLevel := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Severity level",
			Name:        "severity_level",
			Type:        "select",
			Placeholder: "Set the severity level",
			Optional:    true,
		},
		Options:      getDialogOptionsWithSeverityLevels(),
		OptionGroups: []slack.DialogOptionGroup{},
	}

	shouldCreateMeeting := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Should I create an Incident Meeting?",
			Name:        "create_meeting",
			Type:        "select",
			Placeholder: "Select an option",
			Optional:    false,
		},
		Options: []slack.DialogSelectOption{
			{
				Label: "Yes",
				Value: "yes",
			},
			{
				Label: "No",
				Value: "no",
			},
		},
		OptionGroups: []slack.DialogOptionGroup{},
		Value:        "yes",
	}

	description := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident description",
			Name:        "incident_description",
			Type:        "textarea",
			Placeholder: "Brief description on what is happening in this incident. eg. We're having issues with the Product X or Service Y",
			Optional:    true,
		},
		MaxLength: 500,
	}

	dialog := slack.Dialog{
		CallbackID:     "inc-open",
		Title:          "Start an Incident",
		SubmitLabel:    "Start",
		NotifyOnCancel: false,
		Elements: []slack.DialogElement{
			incidentTitle,
			product,
			commander,
			severityLevel,
			shouldCreateMeeting,
			description,
		},
	}

	return app.Client.OpenDialog(triggerID, dialog)
}

// StartIncidentByDialog starts an incident after receiving data from a Slack dialog
func StartIncidentByDialog(
	ctx context.Context,
	app *app.App,
	incidentDetails bot.DialogSubmission,
) error {
	app.Logger.Debug(
		ctx,
		"command/open.StartIncidentByDialog",
		log.NewValue("incident_open_details", incidentDetails),
	)

	var (
		now              = time.Now().UTC()
		incidentAuthor   = incidentDetails.User.ID
		submission       = incidentDetails.Submission
		incidentTitle    = submission["incident_title"]
		product          = submission["product"]
		createMeeting    = submission["create_meeting"]
		commander        = submission["incident_commander"]
		severityLevel    = submission["severity_level"]
		description      = submission["incident_description"]
		environment      = config.Env.Environment
		supportTeam      = config.Env.SupportTeam
		productChannelID = config.Env.ProductChannelID
		meetingURL       = ""
	)

	channelName, err := getChannelNameFromIncidentTitle(incidentTitle)
	if err != nil {
		return err
	}

	user, err := getSlackUserInfo(ctx, app, commander)
	if err != nil {
		return fmt.Errorf("commands.StartIncidentByDialog.get_slack_user_info: incident=%v commanderId=%v error=%v", channelName, commander, err)
	}

	channel, err := app.Client.CreateConversationContext(ctx, channelName, false)
	if err != nil {
		return fmt.Errorf("commands.StartIncidentByDialog.create_conversation_context: incident=%v error=%v", channelName, err)
	}

	severityLevelInt64 := int64(-1)
	if severityLevel != "" {
		severityLevelInt64, err = getStringInt64(severityLevel)
		if err != nil {
			return err
		}
	}

	if createMeeting == "yes" {
		options := map[string]string{
			"channel":     channelName,
			"environment": environment,
		}

		url, err := meeting.CreateMeeting(options)

		if err != nil {
			app.Logger.Error(
				ctx,
				log.Trace(),
				log.Reason("CreateMeetingURL"),
				log.NewValue("error", err),
			)
		}

		meetingURL = url
	}

	incident := model.Incident{
		ChannelName:             channelName,
		ChannelID:               channel.ID,
		Title:                   incidentTitle,
		Product:                 product,
		DescriptionStarted:      description,
		Status:                  model.StatusOpen,
		IdentificationTimestamp: &now,
		SeverityLevel:           severityLevelInt64,
		IncidentAuthor:          incidentAuthor,
		CommanderID:             user.SlackID,
		CommanderEmail:          user.Email,
		MeetingURL:              meetingURL,
	}

	incidentID, err := app.IncidentRepository.InsertIncident(ctx, &incident)
	if err != nil {
		return err
	}

	attachment := createOpenAttachment(incident, incidentID, meetingURL, supportTeam)
	message := "An Incident has been opened by <@" + incident.IncidentAuthor + ">"

	var waitgroup sync.WaitGroup
	defer waitgroup.Wait()

	concurrence.WithWaitGroup(&waitgroup, func() {
		postAndPinMessage(app, channel.ID, message, attachment)
	})
	concurrence.WithWaitGroup(&waitgroup, func() {
		postAndPinMessage(app, productChannelID, message, attachment)
	})

	shouldWritePostMortem := app.FileStorage != nil
	if shouldWritePostMortem {
		//We need run that without wait because the modal need close in only 3s
		go createPostMortemAndFillTopic(ctx, app, incident, incidentID, channel, meetingURL)
	} else {
		fillTopic(ctx, app, incident, channel.ID, meetingURL, "")
	}

	// startReminderStatusJob(ctx, logger, client, repository, incident)

	_, warning, metaWarning, err := app.Client.JoinConversationContext(ctx, channel.ID)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("JoinConversationContext"),
			log.NewValue("warning", warning),
			log.NewValue("meta_warning", metaWarning),
			log.NewValue("error", err),
		)
		return err
	}

	strategy, err := app.Inviter.CreateStrategy(config.Env.InvitationStrategy)
	if err != nil {
		app.Logger.Error(
			ctx,
			"Could not find strategy",
			log.NewValue("StrategyName", config.Env.InvitationStrategy),
			log.NewValue("Error", err),
		)
		return err
	}

	return app.Inviter.InviteStakeholders(ctx, incident, strategy)
}

func createPostMortemAndFillTopic(
	ctx context.Context, app *app.App, incident model.Incident, incidentID int64, channel *slack.Channel, meetingURL string,
) {
	postMortemURL, err := createPostMortem(ctx, app, incidentID, incident.Title, channel.Name)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("createPostMortem"),
			log.NewValue("channel.Name", channel.Name),
			log.NewValue("error", err),
		)
		return
	}

	fillTopic(ctx, app, incident, channel.ID, meetingURL, postMortemURL)
}

func createOpenAttachment(incident model.Incident, incidentID int64, meetingURL string, supportTeam string) slack.Attachment {
	var messageText strings.Builder
	messageText.WriteString("An Incident has been opened by <@" + incident.IncidentAuthor + ">\n\n")
	messageText.WriteString("*Title:* " + incident.Title + "\n")
	messageText.WriteString("*Severity:* " + getSeverityLevelText(incident.SeverityLevel) + "\n\n")
	messageText.WriteString("*Product / Service:* " + incident.Product + "\n")
	messageText.WriteString("*Channel:* <#" + incident.ChannelName + ">\n")
	messageText.WriteString("*Commander:* <@" + incident.CommanderID + ">\n\n")
	messageText.WriteString("*Description:* `" + incident.DescriptionStarted + "`\n\n")
	messageText.WriteString("*Meeting:* " + meetingURL + "\n")

	if supportTeam != "" {
		messageText.WriteString("*cc:* <@" + supportTeam + ">\n")
	}

	preText := ""

	if supportTeam != "" {
		preText = "*cc:* <!subteam^" + supportTeam + ">"
	}

	return slack.Attachment{
		Pretext:  preText,
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#FE4D4D",
		Fields: []slack.AttachmentField{
			{
				Title: "Incident ID",
				Value: strconv.FormatInt(incidentID, 10),
			},
			{
				Title: "Incident Channel",
				Value: "<#" + incident.ChannelID + ">",
			},
			{
				Title: "Incident Title",
				Value: incident.Title,
			},
			{
				Title: "Severity",
				Value: getSeverityLevelText(incident.SeverityLevel),
			},
			{
				Title: "Product / Service",
				Value: incident.Product,
			},
			{
				Title: "Commander",
				Value: "<@" + incident.CommanderID + ">",
			},
			{
				Title: "Description",
				Value: "```" + incident.DescriptionStarted + "```",
			},
			{
				Title: "Meeting",
				Value: meetingURL,
			},
		},
	}
}
