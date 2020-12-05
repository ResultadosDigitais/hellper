package commands

import (
	"context"
	"time"

	"hellper/internal/app"
	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/job"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack"
)

var jobs []job.Job

func canStopReminder(incident model.Incident) bool {
	return incident.Status == model.StatusClosed || incident.Status == model.StatusCancel
}

func requestStatus(ctx context.Context, app *app.App, jobIncident model.Incident) func(j job.Job) {
	return func(j job.Job) {
		incident, err := app.IncidentRepository.GetIncident(ctx, jobIncident.ChannelID)
		if err != nil {
			app.Logger.Error(
				ctx,
				log.Trace(),
				log.Reason("GetIncident"),
				log.NewValue("error", err),
			)
			return
		}

		logWriter := app.Logger.With(
			log.NewValue("channelID", incident.ChannelID),
			log.NewValue("channelName", incident.ChannelName),
		)

		logWriter.Debug(
			ctx,
			log.Trace(),
			log.Action("running"),
		)

		if canStopReminder(incident) {
			logWriter.Debug(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("canStopReminder"),
				log.NewValue("jobIncident.Status", jobIncident.Status),
			)

			job.Stop(&j)
			return
		}

		snoozedUntil := incident.SnoozedUntil
		if snoozedUntil.Time.Unix() > time.Now().Unix() {
			logWriter.Debug(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("isPaused"),
				log.NewValue("snoozedUntil", snoozedUntil.Time),
			)
			return
		}

		if incident.Status != jobIncident.Status {
			logWriter.Debug(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("statusChanged"),
				log.NewValue("jobIncident.Status", jobIncident.Status),
			)
			startReminderStatusJob(ctx, app, incident)
			job.Stop(&j)
			return
		}

		pin, err := bot.LastPin(app.Client, incident.ChannelID)
		if err != nil {
			logWriter.Error(
				ctx,
				log.Trace(),
				log.Reason("LastPin"),
				log.NewValue("error", err),
			)
			return
		}

		if incident.Status == model.StatusResolved {
			now := time.Now()
			endTS := incident.EndTimestamp
			diffHours := now.Sub(*endTS)
			if int(diffHours.Hours()) <= config.Env.SLAHoursToClose {
				logWriter.Debug(
					ctx,
					log.Trace(),
					log.Action("do_not_notify"),
					log.Reason("SLAHoursToClose"),
					log.NewValue("incident.Status", incident.Status),
					log.NewValue("incident.EndTimestamp", incident.EndTimestamp),
					log.NewValue("SLAHoursToClose", config.Env.SLAHoursToClose),
					log.NewValue("diffHours", diffHours),
				)
				return
			}

			sendNotification(ctx, app, incident)
			return
		}

		if pin != (slack.Item{}) {
			timeMessage, err := convertTimestamp(pin.Message.Msg.Timestamp)
			if err != nil {
				logWriter.Error(
					ctx,
					log.Trace(),
					log.Action("convertTimestamp"),
					log.NewValue("error", err),
				)
				return
			}

			if timeMessage.After(time.Now().Add(-setRecurrence(incident))) {
				logWriter.Debug(
					ctx,
					log.Trace(),
					log.Action("do_not_notify"),
					log.Reason("last_pin_time"),
				)
				return
			}
		}

		sendNotification(ctx, app, incident)
	}
}

func startReminderStatusJob(ctx context.Context, app *app.App, incident model.Incident) {
	logWriter := app.Logger.With(
		log.NewValue("channelID", incident.ChannelID),
		log.NewValue("channelName", incident.ChannelName),
		log.NewValue("status", incident.Status),
	)

	logWriter.Debug(
		ctx,
		log.Trace(),
		log.Action("running"),
		log.NewValue("recurrence", setRecurrence(incident).Seconds()),
	)

	j := job.New(setRecurrence(incident), requestStatus(ctx, app, incident))
	jobs = append(jobs, j)
}

// StartAllReminderJobs starts a job for each current active incident. This job posts a reminder in the channel, asking for a incident status update.
// This function is called only once, in the inicialization of the aplication. For new incidents, the startReminderStatusJob is called specifically for that incident.
func StartAllReminderJobs(app *app.App) {
	ctx := context.Background()
	app.Logger.Info(ctx, log.Trace())

	incidents, err := app.IncidentRepository.ListActiveIncidents(ctx)
	if err != nil {
		app.Logger.Error(
			ctx,
			log.Trace(),
			log.Action("ListActiveIncidents"),
			log.NewValue("error", err),
		)
	}

	for _, incident := range incidents {
		startReminderStatusJob(ctx, app, incident)
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

func sendNotification(ctx context.Context, app *app.App, incident model.Incident) {
	logWriter := app.Logger.With(
		log.NewValue("channelID", incident.ChannelID),
		log.NewValue("channelName", incident.ChannelName),
		log.NewValue("incidentStatus", incident.Status),
	)

	logWriter.Info(
		ctx,
		log.Trace(),
		log.Action("postMessage"),
	)

	_, err := postMessage(app, incident.ChannelID, statusNotify(incident))

	if err != nil {
		logWriter.Error(
			ctx,
			log.Trace(),
			log.Action("postMessage"),
			log.NewValue("error", err),
		)
		return
	}
}
