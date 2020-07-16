package commands

import (
	"context"
	"strings"

	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

// OpenCancelIncidentDialog opens a dialog on Slack, so the user can cancel an incident
func OpenCancelIncidentDialog(client bot.Client, triggerID string) error {
	description := &slack.TextInputElement{
		DialogInput: slack.DialogInput{
			Label:       "Description",
			Name:        "incident_description",
			Type:        "textarea",
			Placeholder: "Description eg. We are canceling the suspected incident",
			Optional:    false,
		},
		MaxLength: 500,
	}

	dialog := slack.Dialog{
		CallbackID:     "inc-cancel",
		Title:          "Cancel an Incident",
		SubmitLabel:    "Ok",
		NotifyOnCancel: false,
		Elements: []slack.DialogElement{
			description,
		},
	}

	return client.OpenDialog(triggerID, dialog)
}

// CancelIncidentByDialog cancels an incident after receiving data from a Slack dialog
func CancelIncidentByDialog(ctx context.Context, client bot.Client, logger log.Logger, repository model.Repository, incidentDetails bot.DialogSubmission) error {
	logger.Info(
		ctx,
		"command/cancel.CancelIncidentByDialog",
		log.NewValue("incident_cancel_details", incidentDetails),
	)

	incidentAuthor := incidentDetails.User.ID
	channelID := incidentDetails.Channel.ID
	submission := incidentDetails.Submission
	description := submission.IncidentDescription

	var messageText strings.Builder

	messageText.WriteString("An Incident has been canceled by <@" + incidentAuthor + ">\n\n")
	messageText.WriteString("*Channel:* <#" + channelID + ">\n")
	messageText.WriteString("*Description:* `" + description + "`\n\n")

	attachment := slack.Attachment{
		Pretext:  "",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#EDA248",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Channel",
				Value: "<#" + channelID + ">",
			},
			slack.AttachmentField{
				Title: "Description",
				Value: "```" + description + "```",
			},
		},
	}

	postAndPinMessage(
		client,
		channelID,
		"An Incident has been canceled by <@"+incidentAuthor+"> *cc:* <"+config.Env.SupportTeam+">",
		attachment,
	)
	postAndPinMessage(
		client,
		config.Env.ProductChannelID,
		"An Incident has been canceled by <@"+incidentAuthor+"> *cc:* <"+config.Env.SupportTeam+">",
		attachment,
	)
	repository.CancelIncident(ctx, channelID, description)
	client.ArchiveConversationContext(ctx, channelID)

	return nil
}
