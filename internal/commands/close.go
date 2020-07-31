package commands

import (
	"context"
	"database/sql"
	"hellper/internal/concurrence"
	"strings"
	"sync"

	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

// CloseIncidentDialog opens a dialog on Slack, so the user can close an incident
func CloseIncidentDialog(ctx context.Context, logger log.Logger, client bot.Client, repository model.Repository, channelID, userID, triggerID string) error {
	inc, err := repository.GetIncident(ctx, channelID)
	if err != nil {
		logger.Error(
			ctx,
			"command/dates.UpdateDatesDialog GetIncident ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	if inc.StartTimestamp == nil {
		var (
			messageText strings.Builder
		)

		messageText.WriteString("The dates of Incident <#" + inc.ChannelId + "> has not been updated yet.\n" +
			"Please, call the command `/hellper_update_dates` to receive the current dates and update each one.")

		attch := slack.Attachment{
			Pretext:  "",
			Fallback: messageText.String(),
			Text: "The dates of Incident <#" + inc.ChannelId + "> has not been updated yet.\n" +
				"Please, call the command `/hellper_update_dates` to receive the current dates and update each one.",
			Color:  "#ff8c00",
			Fields: []slack.AttachmentField{},
		}

		return postMessage(client, channelID, "", attch)
	}

	feature := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Feature",
			Name:        "feature",
			Type:        "text",
			Placeholder: "Feature",
			Optional:    false,
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
	responsibility := &slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Label:       "Responsibility",
			Name:        "responsibility",
			Type:        "select",
			Placeholder: "Set the responsible",
			Optional:    false,
		},
		Value: "0",
		Options: []slack.DialogSelectOption{
			{
				Label: "Product",
				Value: "0",
			},
			{
				Label: "Third-Party",
				Value: "1",
			},
		},
		OptionGroups: []slack.DialogOptionGroup{},
	}
	team := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Owner team",
			Name:        "owner_team",
			Type:        "text",
			Placeholder: "Team",
			Optional:    false,
		},
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
	impact := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Impact of incident",
			Name:        "impact",
			Type:        "text",
			Placeholder: "Number of impacted accounts.",
			Optional:    false,
		},
		Subtype: slack.InputSubtypeNumber,
	}

	dialog := slack.Dialog{
		CallbackID:     "inc-close",
		Title:          "Close an Incident",
		SubmitLabel:    "Close",
		NotifyOnCancel: false,
		Elements: []slack.DialogElement{
			impact,
			team,
			feature,
			severityLevel,
			responsibility,
			rootCause,
		},
	}

	return client.OpenDialog(triggerID, dialog)
}

// CloseIncidentByDialog closes an incident after receiving data from a Slack dialog
func CloseIncidentByDialog(ctx context.Context, client bot.Client, logger log.Logger, repository model.Repository, incidentDetails bot.DialogSubmission) error {
	logger.Info(
		ctx,
		"command/close.CloseIncidentByDialog",
		log.NewValue("incident_close_details", incidentDetails),
	)

	var (
		customerImpact   sql.NullInt64
		channelID        = incidentDetails.Channel.ID
		userID           = incidentDetails.User.ID
		userName         = incidentDetails.User.Name
		submissions      = incidentDetails.Submission
		impact           = submissions.Impact
		team             = submissions.Team
		feature          = submissions.Feature
		severityLevel    = submissions.SeverityLevel
		responsibility   = getResponsabilityText(submissions.Responsibility)
		rootCause        = submissions.RootCause
		notifyOnClose    = config.Env.NotifyOnClose
		productChannelID = config.Env.ProductChannelID
		customerImpact   sql.NullInt64
	)

	severityLevelInt64, err := getStringInt64(severityLevel)
	if err != nil {
		return err
	}

	customerImpact.Int64, err = getStringInt64(impact)
	if err != nil {
		return err
	}

	incident := model.Incident{
		RootCause:      rootCause,
		Functionality:  feature,
		Team:           team,
		CustomerImpact: customerImpact,
		SeverityLevel:  severityLevelInt64,
		Responsibility: responsibility,
		ChannelId:      channelID,
	}

	channelAttachment := createCloseChannelAttachment(incident, userName, impact)
	privateAttachment := createClosePrivateAttachment(incident)
	message := "The Incident <#" + incident.ChannelId + "> has been closed by <@" + userName + ">"

	var waitgroup sync.WaitGroup
	defer waitgroup.Wait()

	if notifyOnClose {
		concurrence.WithWaitGroup(&waitgroup, func() {
			postAndPinMessage(
				client,
				productChannelID,
				message,
				channelAttachment,
			)
		})
	}
	concurrence.WithWaitGroup(&waitgroup, func() {
		postMessage(client, userID, "", privateAttachment)
	})

	postAndPinMessage(
		client,
		channelID,
		message,
		channelAttachment,
	)
	err = client.ArchiveConversationContext(ctx, channelID)
	if err != nil {
		logger.Error(
			ctx,
			"command/close.CloseIncidentByDialog ArchiveConversationContext ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("userID", userID),
			log.NewValue("error", err),
		)
		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	repository.CloseIncident(ctx, &incident)

	return nil
}

func getResponsabilityText(r string) string {
	switch r {
	case "0":
		return "Product"
	case "1":
		return "Third-Party"
	}
	return ""
}

func createCloseChannelAttachment(inc model.Incident, userName, impact string) slack.Attachment {
	var messageText strings.Builder
	messageText.WriteString("The Incident <#" + inc.ChannelId + "> has been closed by <@" + userName + ">\n\n")
	messageText.WriteString("*Team:* <#" + inc.Team + ">\n")
	messageText.WriteString("*Feature:* `" + inc.Functionality + "`\n")
	messageText.WriteString("*Impact:* `" + impact + "`\n")
	messageText.WriteString("*Severity:* `" + getSeverityLevelText(inc.SeverityLevel) + "`\n")
	messageText.WriteString("*Responsibility:* `" + inc.Responsibility + "`\n")
	messageText.WriteString("*Root cause:* `" + inc.RootCause + "`\n\n")

	return slack.Attachment{
		Pretext:  "",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#6fff47",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Incident",
				Value: "<#" + inc.ChannelId + ">",
			},
			slack.AttachmentField{
				Title: "Team",
				Value: inc.Team,
			},
			slack.AttachmentField{
				Title: "Feature",
				Value: inc.Functionality,
			},
			slack.AttachmentField{
				Title: "Impact",
				Value: impact,
			},
			slack.AttachmentField{
				Title: "Severity",
				Value: getSeverityLevelText(inc.SeverityLevel),
			},
			slack.AttachmentField{
				Title: "Responsibility",
				Value: inc.Responsibility,
			},
			slack.AttachmentField{
				Title: "RootCause",
				Value: inc.RootCause,
			},
		},
	}
}

func createClosePrivateAttachment(inc model.Incident) slack.Attachment {
	var privateText strings.Builder
	privateText.WriteString("The Incident <#" + inc.ChannelId + "> has been resolved by you\n\n")
	privateText.WriteString("*Status.io:* Be sure to close the incident on https://status.io\n\n")

	return slack.Attachment{
		Pretext:  "The Incident <#" + inc.ChannelId + "> has been closed by you",
		Fallback: privateText.String(),
		Text:     "",
		Color:    "#FE4D4D",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Statu.io",
				Value: "Be sure to close the incident on status.io",
			},
		},
	}
}
