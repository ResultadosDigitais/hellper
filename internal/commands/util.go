package commands

import (
	"context"
	"errors"
	"fmt"
	"hellper/internal/model"
	"regexp"
	"strconv"
	"strings"
	"time"

	"hellper/internal/app"
	"hellper/internal/log"

	"github.com/slack-go/slack"
)

const dateLayout = "2006-01-02T15:04:05-0700"

func ping(ctx context.Context, app *app.App, channelID string) {

	logWriter := app.Logger.With(
		log.NewValue("channelID", channelID),
	)

	_, err := postMessage(app, channelID, "pong")
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

	_, err := postMessage(app, channelID, `
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

func postMessageVisibleOnlyForUser(
	ctx context.Context, app *app.App, channel string, userID string, text string, attachments ...slack.Attachment,
) error {
	_, err := app.Client.PostEphemeralContext(
		ctx,
		channel,
		userID,
		slack.MsgOptionText(text, false),
		slack.MsgOptionAttachments(attachments...),
	)

	return err
}

func postAndPinMessage(app *app.App, channel string, text string, attachment ...slack.Attachment) error {
	msgRef, err := postMessage(app, channel, text, attachment...)
	if err != nil {
		return err
	}

	return pinMessage(app, channel, *msgRef)
}

func postMessage(app *app.App, channelID string, text string, attachments ...slack.Attachment) (*slack.ItemRef, error) {
	return postGenericMessage(
		app,
		channelID,
		text,
		slack.MsgOptionText(text, false),
		slack.MsgOptionAttachments(attachments...),
	)
}

func postGenericMessage(app *app.App, channel string, text string, msgOptions ...slack.MsgOption) (*slack.ItemRef, error) {
	channelID, timestamp, err := app.Client.PostMessage(channel, msgOptions...)
	if err != nil {
		return nil, err
	}

	msgRef := slack.NewRefToMessage(channelID, timestamp)

	return &msgRef, nil
}

func pinMessage(app *app.App, channelID string, msgRef slack.ItemRef) error {
	return app.Client.AddPin(channelID, msgRef)
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

func fillDialogOptionsIfNeeded(options []slack.DialogSelectOption) []slack.DialogSelectOption {
	const minimumNumberOfOptions = 4

	for len(options) < minimumNumberOfOptions {
		options = append(options, slack.DialogSelectOption{
			Label: "----------------",
			Value: "0",
		})
	}

	return options
}

func getDialogOptionsWithServiceInstances(services []*model.ServiceInstance) []slack.DialogSelectOption {
	serviceList := []slack.DialogSelectOption{}

	for _, service := range services {
		serviceList = append(serviceList, slack.DialogSelectOption{
			Label: service.Name,
			Value: service.Name,
		})
	}

	// Slack asks for at least 4 entries in the option panel. So I populate dumby options here, otherwise
	// the open command will fail and will give no feedback whatsoever for the user.
	return fillDialogOptionsIfNeeded(serviceList)
}

func getDialogOptionsWithSeverityLevels() []slack.DialogSelectOption {
	return []slack.DialogSelectOption{
		{
			Label: getSeverityLevelText(0),
			Value: "0",
		},
		{
			Label: getSeverityLevelText(1),
			Value: "1",
		},
		{
			Label: getSeverityLevelText(2),
			Value: "2",
		},
		{
			Label: getSeverityLevelText(3),
			Value: "3",
		},
	}
}

func getChannelNameFromIncidentTitle(incidentTitle string) (string, error) {
	const titleMaxSize = 64

	// first allow only alphanumeric characters on title, based on https://golangcode.com/how-to-remove-all-non-alphanumerical-characters-from-a-string/
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		return "", err
	}

	processedIncidentTitle := strings.ToLower(reg.ReplaceAllString(incidentTitle, ""))

	// then truncate if needed, because Slack supports channel names with an maximum of 80 characters
	if len(processedIncidentTitle) > titleMaxSize { // timeMaxSize is the maximum value (80) excluding the prefix (4, "inc-") and suffix (9, "-yyyyMMdd") to be added
		processedIncidentTitle = processedIncidentTitle[:titleMaxSize]
	}

	// finally, concatenate "inc-" as prefix and a date string as suffix
	currentDate := time.Now()
	currentDateAsString := fmt.Sprintf("%04d%02d%02d", currentDate.Year(), currentDate.Month(), currentDate.Day())
	processedIncidentTitle = fmt.Sprintf("inc-%s-%s", processedIncidentTitle, currentDateAsString)

	return processedIncidentTitle, nil
}

func fillTopic(
	ctx context.Context, app *app.App, incident model.Incident,
	channelID string, meetingURL string, postMortemURL string,
) {
	var topic strings.Builder
	if meetingURL != "" {
		topic.WriteString("*Meeting:* " + meetingURL + "\n\n")
	}
	if postMortemURL != "" {
		topic.WriteString("*PostMortemURL:* " + postMortemURL + "\n\n")
	}
	topic.WriteString("*Commander:* <@" + incident.CommanderID + ">\n\n")
	topicString := topic.String()

	_, err := app.Client.SetTopicOfConversation(channelID, topicString)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("SetTopicOfConversation"),
			log.NewValue("channel.ID", channelID),
			log.NewValue("topic.String", topicString),
			log.NewValue("error", err),
		)
	}
}
