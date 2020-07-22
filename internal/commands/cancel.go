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
func OpenCancelIncidentDialog(
	ctx context.Context,
	logger log.Logger,
	client bot.Client,
	repository model.Repository,
	channelID string,
	userID string,
	triggerID string,
) error {

	inc, err := repository.GetIncident(ctx, channelID)
	if err != nil {
		logger.Error(
			ctx,
			"command/dates.OpenCancelIncidentDialog GetIncident ERROR",
			log.NewValue("channelID", channelID),
			log.NewValue("userID", userID),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	if inc.Status != model.StatusOpen {
		message := "The incident <#" + inc.ChannelId + "> is already `" + inc.Status + "`.\n" +
			"Only a `open` incident can be canceled."

		var messageText strings.Builder
		messageText.WriteString(message)

		attch := slack.Attachment{
			Pretext:  "",
			Fallback: messageText.String(),
			Text:     message,
			Color:    "#ff8c00",
			Fields:   []slack.AttachmentField{},
		}

		_, err = client.PostEphemeralContext(ctx, channelID, userID, slack.MsgOptionAttachments(attch))
		if err != nil {
			logger.Error(
				ctx,
				"command/dates.OpenCancelIncidentDialog PostEphemeralContext ERROR",
				log.NewValue("channelID", channelID),
				log.NewValue("userID", userID),
				log.NewValue("error", err),
			)

			PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
			return err
		}

		return nil
	}

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

	message := "An Incident has been canceled by <@" + incidentAuthor + "> *cc:* <" + config.Env.SupportTeam + ">"

	postAndPinMessage(
		client,
		channelID,
		message,
		attachment,
	)
	postAndPinMessage(
		client,
		config.Env.ProductChannelID,
		message,
		attachment,
	)
	repository.CancelIncident(ctx, channelID, description)
	client.ArchiveConversationContext(ctx, channelID)

	return nil
}
