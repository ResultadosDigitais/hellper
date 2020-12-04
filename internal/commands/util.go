package commands

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"hellper/internal/app"
	"hellper/internal/log"

	"github.com/slack-go/slack"
)

func ping(ctx context.Context, app *app.App, channelID string) {

	logWriter := app.Logger.With(
		log.NewValue("channelID", channelID),
	)

	err := postMessage(app, channelID, "pong")
	if err != nil {
		logWriter.Error(
			ctx,
			"command/util.ping postMessage error",
			log.NewValue("error", err),
		)
	}
}

func help(ctx context.Context, app *app.App, channelID string) {

	logWriter := app.Logger.With(
		log.NewValue("channelID", channelID),
	)

	err := postMessage(app, channelID, `
	hellper
	A bot to help the incident treatment
	Available commands:
 	help      Show this help
 	ping      Test bot connectivity
 	list      List all active incidents
 	state     Show incident state and timeline
`)
	if err != nil {
		logWriter.Error(
			ctx,
			"command/util.help postMessage error",
			log.NewValue("error", err),
		)
	}
}

func getStringInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func getSeverityLevelText(severityLevel int64) string {
	switch severityLevel {
	case 0:
		return "SEV0 - All hands on deck"
	case 1:
		return "SEV1 - Critical impact to many users"
	case 2:
		return "SEV2 - Minor issue that impacts ability to use product"
	case 3:
		return "SEV3 - Minor issue not impacting ability to use product"
	default:
		return ""
	}
}

// PostErrorAttachment posts a error attachment on the given channel
func PostErrorAttachment(ctx context.Context, app *app.App, channelID string, userID string, text string) {
	attach := slack.Attachment{
		Pretext:  "",
		Fallback: "",
		Text:     "",
		Color:    "#FE4D4D",
		Fields: []slack.AttachmentField{
			{
				Title: "Error",
				Value: text,
			},
		},
	}

	_, err := app.Client.PostEphemeralContext(ctx, channelID, userID, slack.MsgOptionAttachments(attach))
	if err != nil {
		app.Logger.Error(
			ctx,
			"command/util.PostErrorAttachment postMessage error",
			log.NewValue("channelID", channelID),
			log.NewValue("userID", userID),
			log.NewValue("text", text),
			log.NewValue("error", err),
		)
		return
	}
}

// PostInfoAttachment posts a info attachment on the given channel
func PostInfoAttachment(ctx context.Context, app *app.App, channelID string, userID string, title string, message string) {
	app.Client.PostEphemeralContext(ctx, channelID, userID, slack.MsgOptionAttachments(slack.Attachment{
		Pretext:  "",
		Fallback: "",
		Text:     "",
		Color:    "#4DA6FE",
		Fields: []slack.AttachmentField{
			{
				Title: title,
				Value: message,
			},
		},
	}))
}

func postMessage(app *app.App, channel string, text string, attachments ...slack.Attachment) error {
	_, _, err := app.Client.PostMessage(channel, slack.MsgOptionText(text, false), slack.MsgOptionAttachments(attachments...))
	if err != nil {
		return err
	}
	return nil
}

func postAndPinMessage(app *app.App, channel string, text string, attachment ...slack.Attachment) error {
	channelID, timestamp, postErr := app.Client.PostMessage(channel, slack.MsgOptionText(text, false), slack.MsgOptionAttachments(attachment...))

	// Grab a reference to the message.
	msgRef := slack.NewRefToMessage(channelID, timestamp)

	// Add message pin to channel
	if postErr = app.Client.AddPin(channelID, msgRef); postErr != nil {
		return postErr
	}
	return nil
}

func convertTimestamp(timestamp string) (time.Time, error) {
	if timestamp == "" {
		return time.Time{}, errors.New("Empty Timestamp")
	}

	timeString := strings.Split(timestamp, ".")
	timeMinutes, err := strconv.ParseInt(timeString[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	timeSec, err := strconv.ParseInt(timeString[1], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	fullTime := time.Unix(timeMinutes, timeSec)

	return fullTime, nil
}
