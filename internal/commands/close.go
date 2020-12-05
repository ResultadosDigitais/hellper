package commands

import (
	"context"
	"hellper/internal/app"
	"hellper/internal/concurrence"
	"strconv"
	"strings"
	"sync"
	"time"

	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

// CloseIncidentDialog opens a dialog on Slack, so the user can close an incident
func CloseIncidentDialog(ctx context.Context, app *app.App, channelID, userID, triggerID string) error {
	inc, err := app.IncidentRepository.GetIncident(ctx, channelID)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("GetIncident"),
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, app, channelID, userID, err.Error())
		return err
	}

	rootCause := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Root Cause",
			Name:        "root_cause",
			Type:        "textarea",
			Placeholder: "Incident root cause description.",
			Optional:    false,
		},
		MaxLength: 500,
	}
	startDate := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Incident start date (UTC)",
			Name:        "init_date",
			Type:        "text",
			Placeholder: dateLayout,
			Hint:        "The time is in format " + dateLayout + " and UTC timezone",
			Optional:    false,
		},
		Value: "",
	}
	severityLevel := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Severity level",
			Name:        "severity_level",
			Type:        "select",
			Placeholder: "Set the severity level",
			Optional:    true,
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

	dialogElements := []slack.DialogElement{
		rootCause,
	}
	if inc.StartTimestamp == nil {
		dialogElements = append(dialogElements, startDate)
	}
	dialogElements = append(dialogElements, severityLevel)

	dialog := slack.Dialog{
		CallbackID:     "inc-close",
		Title:          "Close an Incident",
		SubmitLabel:    "Close",
		NotifyOnCancel: false,
		Elements:       dialogElements,
	}

	return app.Client.OpenDialog(triggerID, dialog)
}

// CloseIncidentByDialog closes an incident after receiving data from a Slack dialog
func CloseIncidentByDialog(ctx context.Context, app *app.App, incidentDetails bot.DialogSubmission) error {
	app.Logger.Debug(
		ctx,
		"command/close.CloseIncidentByDialog",
		log.NewValue("incident_close_details", incidentDetails),
	)

	var (
		channelID        = incidentDetails.Channel.ID
		userID           = incidentDetails.User.ID
		userName         = incidentDetails.User.Name
		submissions      = incidentDetails.Submission
		rootCause        = submissions.RootCause
		startDateText    = submissions.InitDate
		severityLevel    = submissions.SeverityLevel
		notifyOnClose    = config.Env.NotifyOnClose
		productChannelID = config.Env.ProductChannelID

		startDate time.Time
	)

	logWriter := app.Logger.With(
		log.NewValue("channelID", channelID),
		log.NewValue("userID", userID),
	)

	var err error
	if startDateText != "" {
		startDate, err = time.ParseInLocation(dateLayout, startDateText, time.UTC)
		if err != nil {
			logWriter.Error(
				ctx,
				"command/close.CloseIncidentByDialog ParseInLocation start date ERROR",
				log.NewValue("timeZoneString", "UTC"),
				log.NewValue("startDateText", startDateText),
				log.NewValue("error", err),
			)
			PostErrorAttachment(ctx, app, channelID, userID, err.Error())
			return err
		}
	}

	severityLevelInt64 := int64(-1)
	if severityLevel != "" {
		severityLevelInt64, err = getStringInt64(severityLevel)
		if err != nil {
			return err
		}
	}

	inc, err := app.IncidentRepository.GetIncident(ctx, channelID)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("GetIncident"),
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)
		PostErrorAttachment(ctx, app, channelID, userID, err.Error())
		return err
	}

	ownerTeamName, err := app.ServiceRepository.GetServiceInstanceOwnerTeamName(ctx, inc.Product)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("GetServiceInstanceOwnerTeamName"),
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)
		return err
	}

	incident := model.Incident{
		RootCause:      rootCause,
		StartTimestamp: &startDate,
		Team:           ownerTeamName,
		SeverityLevel:  severityLevelInt64,
		ChannelId:      channelID,
	}
	if startDateText != "" {
		incident.StartTimestamp = &startDate
	}

	err = app.IncidentRepository.CloseIncident(ctx, &incident)
	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("CloseIncident"),
			log.NewValue("incident", incident),
			log.NewValue("error", err),
		)
		return err
	}

	inc, err = app.IncidentRepository.GetIncident(ctx, channelID)
	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("GetIncident"),
			log.NewValue("error", err),
		)
		return err
	}

	channelAttachment := createCloseChannelAttachment(inc, userName)
	privateAttachment := createClosePrivateAttachment(inc)
	message := "The Incident <#" + inc.ChannelId + "> has been closed by <@" + userName + ">"

	var waitgroup sync.WaitGroup
	defer waitgroup.Wait()

	if notifyOnClose {
		concurrence.WithWaitGroup(&waitgroup, func() {
			postAndPinMessage(
				app,
				productChannelID,
				message,
				channelAttachment,
			)
		})
	}
	concurrence.WithWaitGroup(&waitgroup, func() {
		postMessage(app, userID, "", privateAttachment)
	})

	postAndPinMessage(
		app,
		channelID,
		message,
		channelAttachment,
	)
	err = app.Client.ArchiveConversationContext(ctx, channelID)
	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Reason("ArchiveConversationContext"),
			log.NewValue("error", err),
		)
		PostErrorAttachment(ctx, app, channelID, userID, err.Error())
		return err
	}

	return nil
}

func createCloseChannelAttachment(inc model.Incident, userName string) slack.Attachment {
	var messageText strings.Builder
	messageText.WriteString("The Incident <#" + inc.ChannelId + "> has been closed by <@" + userName + ">\n\n")
	messageText.WriteString("*Team:* <#" + inc.Team + ">\n")
	messageText.WriteString("*Severity:* `" + getSeverityLevelText(inc.SeverityLevel) + "`\n")
	messageText.WriteString("*Root cause:* `" + inc.RootCause + "`\n\n")

	return slack.Attachment{
		Pretext:  "",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#6fff47",
		Fields: []slack.AttachmentField{
			{
				Title: "Incident ID",
				Value: strconv.FormatInt(inc.Id, 10),
			},
			{
				Title: "Incident Channel",
				Value: "<#" + inc.ChannelId + ">",
			},
			{
				Title: "Incident Title",
				Value: inc.Title,
			},
			{
				Title: "Team",
				Value: inc.Team,
			},
			{
				Title: "Severity",
				Value: getSeverityLevelText(inc.SeverityLevel),
			},
			{
				Title: "RootCause",
				Value: inc.RootCause,
			},
		},
	}
}

func createClosePrivateAttachment(inc model.Incident) slack.Attachment {
	var privateText strings.Builder
	privateText.WriteString("The Incident <#" + inc.ChannelId + "> has been closed by you\n\n")

	return slack.Attachment{
		Pretext:  "The Incident <#" + inc.ChannelId + "> has been closed by you",
		Fallback: privateText.String(),
		Text:     "",
		Color:    "#FE4D4D",
	}
}
