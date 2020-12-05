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

const minimumNumberServices = 4

// OpenStartIncidentDialog opens a dialog on Slack, so the user can start an incident
func OpenStartIncidentDialog(ctx context.Context, app *app.App, triggerID string) error {
	serviceList := []slack.DialogSelectOption{}

	services, err := app.ServiceRepository.ListServiceInstances(ctx)
	if err != nil {
		return err
	}

	// Slack asks for at least 4 entries in the option panel. So I populate dumby options here, otherwise
	// the open command will fail and will give no feedback whatsoever for the user.
	for len(services) < minimumNumberServices {
		services = append(services, &model.ServiceInstance{
			ID:   "-",
			Name: "-",
		})
	}

	for _, service := range services {
		serviceList = append(serviceList, slack.DialogSelectOption{
			Label: service.Name,
			Value: service.Name,
		})
	}

	incidentTitle := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident Title",
			Name:        "incident_title",
			Type:        "text",
			Placeholder: "My Incident Title",
		},
		MaxLength: 100,
	}

	channelName := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Channel name",
			Name:        "channel_name",
			Type:        "text",
			Placeholder: "inc-my-incident",
		},
		MaxLength: 30,
	}

	meeting := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident Room URL",
			Name:        "incident_room_url",
			Type:        "text",
			Placeholder: "Incident Room URL eg. Zoom Join URL, Google Meet URL",
			Optional:    true,
		},
	}

	severityLevel := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Severity level",
			Name:        "severity_level",
			Type:        "select",
			Placeholder: "Set the severity level",
			Optional:    false,
		},
		Options: []slack.DialogSelectOption{
			{
				Label: "SEV0 - All hands on deck",
				Value: "0",
			},
			{
				Label: "SEV1 - Critical impact to many users",
				Value: "1",
			},
			{
				Label: "SEV2 - Minor issue that impacts ability to use product",
				Value: "2",
			},
			{
				Label: "SEV3 - Minor issue not impacting ability to use product",
				Value: "3",
			},
		},
		OptionGroups: []slack.DialogOptionGroup{},
	}

	product := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Product",
			Name:        "product",
			Type:        "select",
			Placeholder: "Set the product",
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
		DataSource:   "users",
		OptionGroups: []slack.DialogOptionGroup{},
	}

	description := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident description",
			Name:        "incident_description",
			Type:        "textarea",
			Placeholder: "Incident description eg. We're having issues with the Product X or Service Y",
			Optional:    false,
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
			channelName,
			meeting,
			severityLevel,
			product,
			commander,
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
		incidentTitle    = submission.IncidentTitle
		channelName      = submission.ChannelName
		incidentRoomURL  = submission.IncidentRoomURL
		severityLevel    = submission.SeverityLevel
		product          = submission.Product
		commander        = submission.IncidentCommander
		description      = submission.IncidentDescription
		environment      = config.Env.Environment
		supportTeam      = config.Env.SupportTeam
		productChannelID = config.Env.ProductChannelID
	)

	user, err := getSlackUserInfo(ctx, app, commander)
	if err != nil {
		return fmt.Errorf("commands.StartIncidentByDialog.get_slack_user_info: incident=%v commanderId=%v error=%v", channelName, commander, err)
	}

	channel, err := app.Client.CreateConversationContext(ctx, channelName, false)
	if err != nil {
		return fmt.Errorf("commands.StartIncidentByDialog.create_conversation_context: incident=%v error=%v", channelName, err)
	}

	severityLevelInt64, err := getStringInt64(severityLevel)
	if err != nil {
		return err
	}

	incident := model.Incident{
		ChannelName:             channelName,
		ChannelId:               channel.ID,
		Title:                   incidentTitle,
		Product:                 product,
		DescriptionStarted:      description,
		Status:                  model.StatusOpen,
		IdentificationTimestamp: &now,
		SeverityLevel:           severityLevelInt64,
		IncidentAuthor:          incidentAuthor,
		CommanderId:             user.SlackID,
		CommanderEmail:          user.Email,
	}

	incidentID, err := app.IncidentRepository.InsertIncident(ctx, &incident)
	if err != nil {
		return err
	}

	if incidentRoomURL == "" {
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

		incidentRoomURL = url
	}

	attachment := createOpenAttachment(incident, incidentID, incidentRoomURL, supportTeam)
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
		go createPostMortemAndFillTopic(ctx, app, incident, incidentID, channel, incidentRoomURL)
	} else {
		fillTopic(ctx, app, incident, channel, incidentRoomURL, "")
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

	_, err = app.Client.InviteUsersToConversationContext(ctx, channel.ID, commander)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("InviteUsersToConversationContext"),
			log.NewValue("channel.ID", channel.ID),
			log.NewValue("commander", commander),
			log.NewValue("error", err),
		)
		return err
	}

	return nil
}

func fillTopic(
	ctx context.Context, app *app.App, incident model.Incident,
	channel *slack.Channel, incidentRoomURL string, postMortemURL string,
) {
	var topic strings.Builder
	if incidentRoomURL != "" {
		topic.WriteString("*IncidentRoom:* " + incidentRoomURL + "\n\n")
	}
	if postMortemURL != "" {
		topic.WriteString("*PostMortemURL:* " + postMortemURL + "\n\n")
	}
	topic.WriteString("*Commander:* <@" + incident.CommanderId + ">\n\n")
	topicString := topic.String()

	_, err := app.Client.SetTopicOfConversation(channel.ID, topicString)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("SetTopicOfConversation"),
			log.NewValue("channel.ID", channel.ID),
			log.NewValue("topic.String", topicString),
			log.NewValue("error", err),
		)
	}
}

func createPostMortemAndFillTopic(
	ctx context.Context, app *app.App, incident model.Incident, incidentID int64, channel *slack.Channel, incidentRoomURL string,
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

	fillTopic(ctx, app, incident, channel, incidentRoomURL, postMortemURL)
}

func createOpenAttachment(incident model.Incident, incidentID int64, incidentRoomURL string, supportTeam string) slack.Attachment {
	var messageText strings.Builder
	messageText.WriteString("An Incident has been opened by <@" + incident.IncidentAuthor + ">\n\n")
	messageText.WriteString("*Title:* " + incident.Title + "\n")
	messageText.WriteString("*Severity:* " + getSeverityLevelText(incident.SeverityLevel) + "\n\n")
	messageText.WriteString("*Product:* " + incident.Product + "\n")
	messageText.WriteString("*Channel:* <#" + incident.ChannelId + ">\n")
	messageText.WriteString("*Commander:* <@" + incident.CommanderId + ">\n\n")
	messageText.WriteString("*Description:* `" + incident.DescriptionStarted + "`\n\n")
	messageText.WriteString("*Incident Room:* " + incidentRoomURL + "\n")
	messageText.WriteString("*cc:* <@" + supportTeam + ">\n")

	return slack.Attachment{
		Pretext:  "*cc:* <!subteam^" + supportTeam + ">",
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
				Value: "<#" + incident.ChannelId + ">",
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
				Title: "Product",
				Value: incident.Product,
			},
			{
				Title: "Commander",
				Value: "<@" + incident.CommanderId + ">",
			},
			{
				Title: "Description",
				Value: "```" + incident.DescriptionStarted + "```",
			},
			{
				Title: "Incident Room",
				Value: incidentRoomURL,
			},
		},
	}
}
