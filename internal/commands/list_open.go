package commands

import (
	"context"
	"hellper/internal/app"
	"hellper/internal/log"
	"hellper/internal/model"
	"strings"

	"github.com/slack-go/slack"
)

//ListOpenIncidents get the currently opened incidents and return the channel_name of each one of them.
func ListOpenIncidents(ctx context.Context, app *app.App, event TriggerEvent) {
	loggerWriter := app.Logger.With(
		log.NewValue("event", event),
	)

	incidents, err := app.IncidentRepository.ListActiveIncidents(ctx)
	if err != nil {
		loggerWriter.Error(
			ctx,
			"command/list_open.ListOpenIncidents ListActiveIncidents ERROR",
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, app, event.Channel, event.User, err.Error())
	}

	loggerWriter.Debug(
		ctx,
		"command/list_open.ListOpenIncidents",
		log.NewValue("incidents", incidents),
	)

	if len(incidents) == 0 {
		var messageText strings.Builder
		messageText.WriteString("No active incidents!")
	} else {
		attachment := createListOpenAttachment(incidents)
		postMessage(app, event.Channel, "", attachment)
	}
}

func createListOpenAttachment(incidents []model.Incident) slack.Attachment {
	var messageText strings.Builder
	messageText.WriteString("Current open incidents:")

	var fields []slack.AttachmentField

	for _, inc := range incidents {
		messageText.WriteString("- <#" + inc.ChannelID + ">\n")

		fields = append(
			fields,
			slack.AttachmentField{
				Value: "- <#" + inc.ChannelID + ">",
			},
		)
	}

	return slack.Attachment{
		Pretext:  "Current open incidents:",
		Fallback: messageText.String(),
		Text:     "",
		Color:    "#000000",
		Fields:   fields,
	}
}
