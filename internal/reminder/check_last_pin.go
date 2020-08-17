package reminder

import (
	"context"
	"errors"
	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func hasLastPin(ctx context.Context, client bot.Client, logger log.Logger, incident model.Incident) bool {
	pin, err := bot.LastPin(client, incident.ChannelId)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("LastPin"),
			log.NewValue("channelID", incident.ChannelId),
			log.NewValue("channelName", incident.ChannelName),
			log.NewValue("error", err),
		)
		return true
	}

	if pin != (slack.Item{}) {
		timeMessage, err := convertTimestamp(pin.Message.Msg.Timestamp)
		if err != nil {
			logger.Error(
				ctx,
				log.Trace(),
				log.Action("convertTimestamp"),
				log.NewValue("channelID", incident.ChannelId),
				log.NewValue("channelName", incident.ChannelName),
				log.NewValue("error", err),
			)
			return true
		}

		if timeMessage.After(time.Now().Add(-setRecurrence(incident))) {
			logger.Info(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("last_pin_time"),
				log.NewValue("channelID", incident.ChannelId),
				log.NewValue("channelName", incident.ChannelName),
			)
			return true
		}
	}

	return false
}

func setRecurrence(incident model.Incident) time.Duration {
	switch incident.Status {
	case model.StatusOpen:
		return time.Duration(config.Env.ReminderOpenStatusSeconds) * time.Second
	case model.StatusResolved:
		return time.Duration(config.Env.ReminderResolvedStatusSeconds) * time.Second
	}
	return 0
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
