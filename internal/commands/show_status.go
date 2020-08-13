package commands

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"time"

	"hellper/internal/bot"
	"hellper/internal/config"
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
			"command/show_status.createDatesAttachment GetIncident error",
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
			"command/show_status.createStatusAttachment ListPins error",
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
					"command/show_status.createStatusAttachment convertTimestamp error",
					log.NewValue("channelID", channelID),
					log.NewValue("error", err),
				)

				return slack.Attachment{}, err
			}

			if item.Message.User != "" {
				logger.Info(
					ctx,
					"TESTE",
					log.NewValue("message", item.Message),
				)
				user, err := client.GetUserInfoContext(ctx, item.Message.User)
				if err != nil {
					logger.Error(
						ctx,
						"command/show_status.createStatusAttachment GetUserInfoContext error",
						log.NewValue("channelID", channelID),
						log.NewValue("error", err),
					)

					return slack.Attachment{}, err
				}

				msg := treatMessage(ctx, client, logger, item.Message.Text)

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
			Title: config.Env.Messages.NoTimelineItems,
		}
		fields = append(fields, field)

		attach = slack.Attachment{
			Pretext:  "Incident Status:",
			Fallback: "Incident Status",
			Text:     "",
			Color:    "#999999",
			Fields:   fields,
		}
	}
	return attach, nil
}

func treatMessage(ctx context.Context, client bot.Client, logger log.Logger, msg string) string {
	re := regexp.MustCompile(`<@(\w+)>`)
	x := re.FindAllStringSubmatch(msg, -1)

	for _, y := range x {
		user, _ := client.GetUserInfoContext(ctx, y[1])
		msg = strings.Replace(msg, y[0], user.Name, -1)
	}

	return msg
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
		"command/show_status.ShowStatus",
		log.NewValue("channelID", channelID),
	)

	attachDates, err := createDatesAttachment(ctx, logger, repository, channelID)
	if err != nil {
		logger.Error(
			ctx,
			"command/show_status.ShowStatus createDatesAttachment error",
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	attachStatus, err = createStatusAttachment(ctx, client, logger, channelID)
	if err != nil {
		logger.Error(
			ctx,
			"command/show_status.ShowStatus createStatusAttachment error",
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		PostErrorAttachment(ctx, client, logger, channelID, userID, err.Error())
		return err
	}

	postMessage(client, channelID, "", attachDates, attachStatus)
	return nil
}
