package commands

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"time"

	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

func createDateFields(inc model.Incident) (fields []slack.AttachmentField) {
	dateLayout := time.RFC1123

	if startTime := inc.StartTimestamp; startTime != nil {
		timeMessage := startTime.Format(dateLayout)

		field := slack.AttachmentField{
			Title: "Incident Initial Time:",
			Value: timeMessage,
		}
		fields = append(fields, field)
	}

	if identificationTime := inc.IdentificationTimestamp; identificationTime != nil {
		timeMessage := identificationTime.Format(dateLayout)

		field := slack.AttachmentField{
			Title: "Incident Identification Time:",
			Value: timeMessage,
		}
		fields = append(fields, field)
	}

	if endTime := inc.EndTimestamp; endTime != nil {
		timeMessage := endTime.Format(dateLayout)

		field := slack.AttachmentField{
			Title: "Incident End Time:",
			Value: timeMessage,
		}
		fields = append(fields, field)
	}

	return fields
}

func createDatesAttachment(ctx context.Context, logger log.Logger, repository model.Repository, channelID string) (slack.Attachment, error) {
	inc, err := repository.GetIncident(ctx, channelID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("GetIncident"),
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		return slack.Attachment{}, err
	}

	fields := createDateFields(inc)
	attach := slack.Attachment{
		Pretext:  "Incident Dates:",
		Fallback: "Incident Dates",
		Text:     "",
		Color:    "#f2b12e",
		Fields:   fields,
	}

	return attach, nil
}

func createStatusAttachment(ctx context.Context, client bot.Client, logger log.Logger, channelID string) (slack.Attachment, error) {
	var (
		attach     slack.Attachment
		fields     []slack.AttachmentField
		attachText string
	)

	items, _, err := client.ListPins(channelID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("ListPins"),
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		return slack.Attachment{}, err
	}

	sort.Slice(
		items,
		func(i, j int) bool {
			return items[i].Message.Timestamp < items[j].Message.Timestamp
		},
	)

	if len(items) > 0 {
		for _, item := range items {
			attachText = ""

			timeMessage, err := convertTimestamp(item.Message.Timestamp)
			if err != nil {
				logger.Error(
					ctx,
					log.Trace(),
					log.Reason("convertTimestamp"),
					log.NewValue("channelID", channelID),
					log.NewValue("error", err),
				)

				return slack.Attachment{}, err
			}

			if item.Message.User != "" {
				user, err := client.GetUserInfoContext(ctx, item.Message.User)
				if err != nil {
					logger.Error(
						ctx,
						log.Trace(),
						log.Reason("GetUserInfoContext"),
						log.NewValue("channelID", channelID),
						log.NewValue("error", err),
					)

					return slack.Attachment{}, err
				}

				msg, err := treatMessage(ctx, client, logger, item.Message.Text)
				if err != nil {
					return slack.Attachment{}, err
				}

				attachText = msg + " - @" + user.Name
			} else {
				attachText = item.Message.Attachments[0].Pretext + " - @Hellper"
			}

			field := slack.AttachmentField{
				Value: "```" +
					timeMessage.Format(time.RFC1123) +
					"\n" +
					attachText +
					"```",
			}
			fields = append(fields, field)
		}

		attach = slack.Attachment{
			Pretext:  "Incident Status:",
			Fallback: "Incident Status",
			Text:     "",
			Color:    "#f2b12e",
			Fields:   fields,
		}
	} else {
		field := slack.AttachmentField{
			Title: "Incident Timeline is empty",
		}
		fields = append(fields, field)

		attach = slack.Attachment{
			Pretext:  "Incident Status:",
			Fallback: "Incident Status",
			Text:     "",
			Color:    "#999999",
			Fields:   fields,
		}
		logger.Info(
			ctx,
			log.Trace(),
			log.Reason("AttachmentField"),
			log.NewValue("channelID", channelID),
		)
	}
	return attach, nil
}

func treatMessage(ctx context.Context, client bot.Client, logger log.Logger, msg string) (string, error) {
	msg = treatHere(msg)
	msg, err := treatUsersMentions(ctx, client, logger, msg)
	if err != nil {
		return "", err
	}

	msg, err = treatGroupMentions(ctx, client, logger, msg)
	if err != nil {
		return "", err
	}

	return msg, nil
}

func treatHere(msg string) string {
	x := []string{
		"here",
		"channel",
	}

	for _, w := range x {
		msg = strings.Replace(msg, "<!"+w+">", "@"+w, -1)
	}

	return msg
}

func treatUsersMentions(ctx context.Context, client bot.Client, logger log.Logger, msg string) (string, error) {
	re := regexp.MustCompile(`<@(\w+)>`)
	userIDs := re.FindAllStringSubmatch(msg, -1)

	for _, id := range userIDs {
		user, err := client.GetUserInfoContext(ctx, id[1])
		if err != nil {
			logger.Error(
				ctx,
				log.Trace(),
				log.Reason("GetUserInfoContext"),
				log.NewValue("message", msg),
				log.NewValue("error", err),
			)
			return "", err
		}

		msg = strings.Replace(msg, id[0], "@"+user.Name, -1)
	}

	return msg, nil
}

func treatGroupMentions(ctx context.Context, client bot.Client, logger log.Logger, msg string) (string, error) {
	re := regexp.MustCompile(`<!subteam\^(\w+[^>]*)>`)
	groupIDs := re.FindAllStringSubmatch(msg, -1)

	logger.Info(
		ctx,
		log.Trace(),
		log.Reason("treatGroupMentions"),
		log.NewValue("message", msg),
		log.NewValue("goupIDs", groupIDs),
		log.NewValue("info", "untreated group mention"),
	)

	for _, id := range groupIDs {
		group, err := client.GetGroupInfoContext(ctx, id[1])
		if err != nil {
			logger.Error(
				ctx,
				log.Trace(),
				log.Reason("GetGroupInfoContext"),
				log.NewValue("message", msg),
				log.NewValue("error", err),
			)
			return "", err
		}

		msg = strings.Replace(msg, id[0], "@"+group.GroupConversation.Name, -1)
	}

	return msg, nil
}

// ShowStatus posts an attachment on the channel, with each pinned message from it
func ShowStatus(
	ctx context.Context,
	client bot.Client,
	logger log.Logger,
	repository model.Repository,
	channelID string,
	userID string,
) error {

	var (
		attachDates  slack.Attachment
		attachStatus slack.Attachment
	)

	logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("channelID", channelID),
	)

	attachDates, err := createDatesAttachment(ctx, logger, repository, channelID)
	if err != nil {
		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	attachStatus, err = createStatusAttachment(ctx, client, logger, channelID)
	if err != nil {
		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	postMessage(client, channelID, "", attachDates, attachStatus)
	return nil
}
