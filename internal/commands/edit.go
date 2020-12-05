package commands

import (
	"context"
	"fmt"
	"hellper/internal/app"
	"hellper/internal/config"
	"hellper/internal/model"
	"strconv"
	"strings"
	"time"

	"hellper/internal/bot"
	"hellper/internal/log"

	"github.com/slack-go/slack"
)

// OpenEditIncidentDialog opens a dialog on Slack, so the user can edit an incident
func OpenEditIncidentDialog(ctx context.Context, app *app.App, channelID string, triggerID string) error {
	var (
		startTimestampAsText = ""
	)

	services, err := app.ServiceRepository.ListServiceInstances(ctx)
	if err != nil {
		return err
	}

	serviceList := getDialogOptionsWithServiceInstances(services)

	inc, err := app.IncidentRepository.GetIncident(ctx, channelID)
	if err != nil {
		return err
	}

	if inc.StartTimestamp != nil {
		startTimestampAsText = inc.StartTimestamp.Format(dateLayout)
	}

	incidentTitle := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident Title",
			Name:        "incident_title",
			Type:        "text",
			Placeholder: "My Incident Title",
		},
		MaxLength: 100,
		Value:     inc.Title,
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
		Value:        inc.Product,
	}

	commander := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Incident commander",
			Name:        "incident_commander",
			Type:        "select",
			Placeholder: "Set the Incident commander",
			Optional:    false,
		},
		Value:        inc.CommanderID,
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
		Value:        fmt.Sprintf("%d", inc.SeverityLevel),
	}

	meeting := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Meeting URL",
			Name:        "meeting_url",
			Type:        "text",
			Placeholder: "Meeting URL used to discuss the incident eg. Zoom Join URL, Google Meet URL",
			Optional:    true,
		},
		Value: inc.MeetingURL,
	}

	postMortem := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "PostMortem URL",
			Name:        "post_mortem_url",
			Type:        "text",
			Placeholder: "PostMortem URL used to discuss and learn about the incident  eg. Google Docs URL, Wiki URL",
			Optional:    true,
		},
		Value: inc.PostMortemURL,
	}

	startDate := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Start date",
			Name:        "init_date",
			Type:        "text",
			Placeholder: dateLayout,
			Optional:    true,
		},
		Hint:  "The time is in format " + dateLayout,
		Value: startTimestampAsText,
	}

	rootCause := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Root Cause",
			Name:        "root_cause",
			Type:        "textarea",
			Placeholder: "Incident root cause description.",
			Optional:    true,
		},
		MaxLength: 500,
		Value:     inc.RootCause,
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

	// Slack force us to have a maximum of 10 fields in the dialog
	dialog := slack.Dialog{
		CallbackID:     "inc-edit",
		Title:          "Edit an Incident",
		SubmitLabel:    "Save",
		NotifyOnCancel: false,
		Elements: []slack.DialogElement{
			incidentTitle,
			product,
			commander,
			severityLevel,
			meeting,
			postMortem,
			startDate,
			rootCause,
			description,
		},
	}

	return app.Client.OpenDialog(triggerID, dialog)
}

// EditIncidentByDialog starts an incident after receiving data from a Slack dialog
func EditIncidentByDialog(
	ctx context.Context,
	app *app.App,
	incidentDetails bot.DialogSubmission,
) error {
	app.Logger.Info(
		ctx,
		"command/open.EditIncidentByDialog",
		log.NewValue("incident_edit_details", incidentDetails),
	)

	var (
		userID             = incidentDetails.User.ID
		channelID          = incidentDetails.Channel.ID
		channelName        = incidentDetails.Channel.Name
		submission         = incidentDetails.Submission
		incidentTitle      = submission["incident_title"]
		product            = submission["product"]
		commander          = submission["incident_commander"]
		severityLevel      = submission["severity_level"]
		meeting            = submission["meeting_url"]
		postMortem         = submission["post_mortem_url"]
		startTimestampText = submission["init_date"]
		rootCause          = submission["root_cause"]
		description        = submission["incident_description"]
		supportTeam        = config.Env.SupportTeam
		startTimestamp     time.Time
	)

	incidentBeforeEdit, err := app.IncidentRepository.GetIncident(ctx, channelID)
	if err != nil {
		return err
	}

	user, err := getSlackUserInfo(ctx, app, commander)
	if err != nil {
		return fmt.Errorf("commands.EditIncidentByDialog.get_slack_user_info: incident=%v commanderId=%v error=%v", channelName, commander, err)
	}

	severityLevelInt64 := int64(-1)
	if severityLevel != "" {
		severityLevelInt64, err = getStringInt64(severityLevel)
		if err != nil {
			return err
		}
	}

	if startTimestampText != "" {
		startTimestamp, err = time.Parse(dateLayout, startTimestampText)
		if err != nil {
			app.Logger.Error(
				ctx,
				"command.EditIncidentByDialog Parse ERROR",
				log.NewValue("channelID", channelID),
				log.NewValue("startTimestampText", startTimestampText),
				log.NewValue("error", err),
			)

			PostErrorAttachment(ctx, app, channelID, userID, err.Error())
			return err
		}

		// convert date to timestamp
		startTimestamp = startTimestamp.UTC()
	}

	incident := model.Incident{
		ID:                 incidentBeforeEdit.ID,
		Title:              incidentTitle,
		Product:            product,
		DescriptionStarted: description,
		SeverityLevel:      severityLevelInt64,
		CommanderID:        user.SlackID,
		CommanderEmail:     user.Email,
		MeetingURL:         meeting,
		PostMortemURL:      postMortem,
		RootCause:          rootCause,
	}

	if startTimestampText != "" {
		incident.StartTimestamp = &startTimestamp
	}

	err = app.IncidentRepository.UpdateIncident(ctx, &incident)

	if err != nil {
		return err
	}

	if incidentBeforeEdit.CommanderID != incident.CommanderID ||
		incidentBeforeEdit.PostMortemURL != incident.PostMortemURL ||
		incidentBeforeEdit.MeetingURL != incident.MeetingURL {
		fillTopic(ctx, app, incident, channelID, meeting, postMortem)
	}

	attachment := createEditAttachment(incident, incidentBeforeEdit.ID, meeting, supportTeam, incidentDetails.User.Name)
	message := fmt.Sprintf("The incident %d has been edited by <@%s>\n\n", incident.ID, incidentDetails.User.Name)

	postAndPinMessage(app, channelID, message, attachment)

	return nil
}

func createEditAttachment(incident model.Incident, incidentID int64, meetingURL string, supportTeam string, editorName string) slack.Attachment {
	var messageText strings.Builder
	messageText.WriteString(fmt.Sprintf("The incident %d has been edited by <@%s>\n\n", incidentID, editorName))
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
