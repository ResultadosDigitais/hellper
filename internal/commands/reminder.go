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

var jobs []job.Job

func canStopReminder(incident model.Incident) bool {
	return incident.Status == model.StatusClosed || incident.Status == model.StatusCancel
}

func requestStatus(ctx context.Context, client bot.Client, logger log.Logger, repository model.Repository, jobIncident model.Incident) func(j job.Job) {
	return func(j job.Job) {
		incident, err := repository.GetIncident(ctx, jobIncident.ChannelId)
		if err != nil {
			logger.Error(
				ctx,
				log.Trace(),
				log.Reason("GetIncident"),
				log.NewValue("channelID", incident.ChannelId),
				log.NewValue("channelName", incident.ChannelName),
				log.NewValue("error", err),
			)
			return
		}

		logger.Info(
			ctx,
			log.Trace(),
			log.Action("running"),
			log.NewValue("channelID", incident.ChannelId),
			log.NewValue("channelName", incident.ChannelName),
		)

		if canStopReminder(incident) {
			logger.Info(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("canStopReminder"),
				log.NewValue("channelID", incident.ChannelId),
				log.NewValue("channelName", incident.ChannelName),
				log.NewValue("incident.Status", incident.Status),
				log.NewValue("jobIncident.Status", jobIncident.Status),
			)

			job.Stop(&j)
			return
		}

		snoozedUntil := incident.SnoozedUntil
		if snoozedUntil.Time.Unix() > time.Now().Unix() {
			logger.Info(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("isPaused"),
				log.NewValue("channelID", incident.ChannelId),
				log.NewValue("channelName", incident.ChannelName),
				log.NewValue("snoozedUntil", snoozedUntil.Time),
			)
			return
		}

		if incident.Status != jobIncident.Status {
			logger.Info(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("statusChanged"),
				log.NewValue("channelID", incident.ChannelId),
				log.NewValue("channelName", incident.ChannelName),
				log.NewValue("incident.Status", incident.Status),
				log.NewValue("jobIncident.Status", jobIncident.Status),
			)
			startReminderStatusJob(ctx, logger, client, repository, incident)
			job.Stop(&j)
			return
		}

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
			return
		}

		if incident.Status == model.StatusResolved {
			now := time.Now()
			endTS := incident.EndTimestamp
			diffHours := now.Sub(*endTS)
			if int(diffHours.Hours()) <= config.Env.SLAHoursToClose {
				logger.Info(
					ctx,
					log.Trace(),
					log.Action("do_not_notify"),
					log.Reason("SLAHoursToClose"),
					log.NewValue("channelID", incident.ChannelId),
					log.NewValue("channelName", incident.ChannelName),
					log.NewValue("incident.Status", incident.Status),
					log.NewValue("incident.EndTimestamp", incident.EndTimestamp),
					log.NewValue("SLAHoursToClose", config.Env.SLAHoursToClose),
					log.NewValue("diffHours", diffHours),
				)
				return
			}

			sendNotification(ctx, logger, client, incident)
			return
		}

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
			return
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
			return
		}

		sendNotification(ctx, logger, client, incident)
	}
}

func startReminderStatusJob(ctx context.Context, logger log.Logger, client bot.Client, repository model.Repository, incident model.Incident) {
	logger.Info(
		ctx,
		log.Trace(),
		log.Action("running"),
		log.NewValue("channelID", incident.ChannelId),
		log.NewValue("ChannelName", incident.ChannelName),
		log.NewValue("Status", incident.Status),
		log.NewValue("recurrence", setRecurrence(incident).Seconds()),
	)

	j := job.New(setRecurrence(incident), requestStatus(ctx, client, logger, repository, incident))
	jobs = append(jobs, j)
}

// StartAllReminderJobs starts a job for each current active incident. This job posts a reminder in the channel, asking for a incident status update.
// This function is called only once, in the inicialization of the aplication. For new incidents, the startReminderStatusJob is called specifically for that incident.
func StartAllReminderJobs(logger log.Logger, client bot.Client, repository model.Repository) {
	ctx := context.Background()
	logger.Info(ctx, log.Trace())

	incidents, err := repository.ListActiveIncidents(ctx)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("ListActiveIncidents"),
			log.NewValue("error", err),
		)
	}

	for _, incident := range incidents {
		startReminderStatusJob(ctx, logger, client, repository, incident)
	}

}

func statusNotify(incident model.Incident) string {
	switch incident.Status {
	case model.StatusOpen:
		return config.Env.ReminderOpenNotifyMsg
	case model.StatusResolved:
		return config.Env.ReminderResolvedNotifyMsg
	}
	return ""
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

func sendNotification(ctx context.Context, logger log.Logger, client bot.Client, incident model.Incident) {
	err := postMessage(client, incident.ChannelId, statusNotify(incident))

	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("postMessage"),
			log.NewValue("channelID", incident.ChannelId),
			log.NewValue("channelName", incident.ChannelName),
			log.NewValue("incident.Status", incident.Status),
			log.NewValue("error", err),
		)
		return
	}

	logger.Info(
		ctx,
		log.Trace(),
		log.Action("postMessage"),
		log.NewValue("channelID", incident.ChannelId),
		log.NewValue("channelName", incident.ChannelName),
		log.NewValue("incident.Status", incident.Status),
	)
}
