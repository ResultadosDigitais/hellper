package commands

import (
	"context"
	"time"

	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/job"
	"hellper/internal/log"
	"hellper/internal/model"
)

var recurrence = time.Duration(config.Env.ReminderStatusSeconds) * time.Second

var jobs []job.Job

func canStopReminder(incident model.Incident) bool {
	return incident.Status == model.StatusResolved ||
		incident.Status == model.StatusClosed ||
		incident.Status == model.StatusCancel
}

func requestStatus(ctx context.Context, client bot.Client, logger log.Logger, repository model.Repository, channelID string) func(j job.Job) {
	return func(j job.Job) {
		incident, err := repository.GetIncident(ctx, channelID)
		if err != nil {
			logger.Error(
				ctx,
				"command/reminder.requestStatus GetIncident error",
				log.NewValue("channelID", channelID),
				log.NewValue("error", err),
			)
			return
		}

		if canStopReminder(incident) {
			logger.Info(
				ctx,
				"command/reminder.requestStatus stop",
				log.NewValue("channelID", channelID),
				log.NewValue("job", j),
			)

			job.Stop(&j)
			return
		}

		logger.Info(
			ctx,
			"command/reminder.requestStatus running",
			log.NewValue("channelID", channelID),
			log.NewValue("job", j),
		)
		pin, err := bot.LastPin(client, channelID)
		if err != nil {
			logger.Error(
				ctx,
				"command/reminder.requestStatus LastPin error",
				log.NewValue("channelID", channelID),
				log.NewValue("error", err),
			)
			return
		}

		timeMessage, err := convertTimestamp(pin.Message.Msg.Timestamp)
		if err != nil {
			logger.Error(
				ctx,
				"command/reminder.requestStatus convertTimestamp error",
				log.NewValue("channelID", channelID),
				log.NewValue("error", err),
			)
			return
		}

		if timeMessage.Before(time.Now().Add(-recurrence)) {
			err := postMessage(client, channelID, "Update the status of this incident, just pin a message with status on the channel.")
			if err != nil {
				logger.Error(
					ctx,
					"command/reminder.requestStatus postMessage error",
					log.NewValue("channelID", channelID),
					log.NewValue("error", err),
				)
			}
		} else {
			logger.Info(
				ctx,
				"command/reminder.requestStatus OK",
				log.NewValue("channelID", channelID),
			)
		}
	}
}

func startReminderStatusJob(ctx context.Context, logger log.Logger, client bot.Client, repository model.Repository, channelID string) {
	logger.Info(
		ctx,
		"command/reminder.startReminderStatusJob",
		log.NewValue("channelID", channelID),
		log.NewValue("recurrence", recurrence.Seconds()),
	)

	j := job.New(recurrence, requestStatus(ctx, client, logger, repository, channelID))
	jobs = append(jobs, j)
}

// StartAllReminderJobs starts a job for each current active incident. This job posts a reminder in the channel, asking for a incident status update.
// This function is called only once, in the inicialization of the aplication. For new incidents, the startReminderStatusJob is called specifically for that incident.
func StartAllReminderJobs(logger log.Logger, client bot.Client, repository model.Repository) {
	ctx := context.Background()

	logger.Info(
		ctx,
		"command/reminder.StartAllReminderJobs",
		log.NewValue("recurrence", recurrence.Seconds()),
	)

	incidents, err := repository.ListActiveIncidents(ctx)
	if err != nil {
		logger.Error(
			ctx,
			"command/reminder.StartAllReminderJobs ListActiveIncidents error",
			log.NewValue("error", err),
		)
	}

	for _, incident := range incidents {
		startReminderStatusJob(ctx, logger, client, repository, incident.ChannelId)
	}
}
