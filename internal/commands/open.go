package commands

import (
	"context"
	"fmt"
	"hellper/internal/app"
	"hellper/internal/concurrence"
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
			Label:       "Incident Commander",
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
			Label:       "Severity Level",
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
			Label:       "Should I Create an Incident Meeting?",
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
			Label:       "Incident Description",
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
		commanderSlackID = submission["incident_commander"]
		severityLevel    = submission["severity_level"]
		description      = submission["incident_description"]
		environment      = config.Env.Environment
		productChannelID = config.Env.ProductChannelID
		meetingURL       = ""
	)

	channelName, err := getChannelNameFromIncidentTitle(incidentTitle)
	if err != nil {
		return err
	}

	commander, err := getSlackUserInfo(ctx, app, commanderSlackID)
	if err != nil {
		return fmt.Errorf("commands.StartIncidentByDialog.get_slack_user_info: incident=%v commanderId=%v error=%v", channelName, commanderSlackID, err)
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
		CommanderID:             commander.SlackID,
		CommanderEmail:          commander.Email,
		MeetingURL:              meetingURL,
	}

	incidentID, err := app.IncidentRepository.InsertIncident(ctx, &incident)
	if err != nil {
		return err
	}

	card := createOpenCard(incident, incidentID, commander)

	var waitgroup sync.WaitGroup
	defer waitgroup.Wait()

	concurrence.WithWaitGroup(&waitgroup, func() {
		postAndPinBlockMessage(app, channel.ID, card)
	})
	concurrence.WithWaitGroup(&waitgroup, func() {
		postAndPinBlockMessage(app, productChannelID, card)
	})

	fillTopic(ctx, app, incident, channel.ID, meetingURL, "")

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

func createTextBlock(text string, opts ...interface{}) *slack.TextBlockObject {
	blockMessage := fmt.Sprintf(text, opts)
	return slack.NewTextBlockObject("mrkdwn", blockMessage, false, false)
}

func createOpenCard(incident model.Incident, incidentID int64, commander *model.User) []slack.Block {
	headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(":warning: *Incident #%d - %s*", incidentID, incident.Title), false, false)
	headerBlock := slack.NewSectionBlock(headerText, nil, nil)

	bodySlice := []string{}

	bodySlice = append(bodySlice, fmt.Sprintf("*Product / Service:*\t%s", incident.Product))
	bodySlice = append(bodySlice, fmt.Sprintf("*Channel:*\t\t\t\t\t#%s", incident.ChannelName))
	bodySlice = append(bodySlice, fmt.Sprintf("*Commander:*\t\t\t@%s", commander.DisplayName))

	if incident.SeverityLevel > 0 {
		bodySlice = append(bodySlice, fmt.Sprintf("*Severity:*\t\t\t\t\t%s", getSeverityLevelText(incident.SeverityLevel)))
	}

	if incident.MeetingURL != "" {
		bodySlice = append(bodySlice, fmt.Sprintf("*Meeting:*\t\t\t\t\t<%s|access meeting room>", incident.MeetingURL))
	}

	if incident.DescriptionStarted != "" {
		bodySlice = append(bodySlice, fmt.Sprintf("\n*Description:*\n%s", incident.DescriptionStarted))
	}

	dividerBlock := slack.NewDividerBlock()

	bodyBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", strings.Join(bodySlice, "\n"), false, false),
		nil,
		nil,
	)

	return []slack.Block{headerBlock, dividerBlock, bodyBlock}
}

func postAndPinBlockMessage(app *app.App, channel string, blockMessage []slack.Block) error {
	channelID, timestamp, err := app.Client.PostMessage(channel, slack.MsgOptionBlocks(blockMessage...))
	if err != nil {
		return err
	}

	msgRef := slack.NewRefToMessage(channelID, timestamp)

	return pinMessage(app, channel, msgRef)
}
