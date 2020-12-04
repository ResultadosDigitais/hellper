package reminder

import (
	"context"
	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"
	"time"
)

func hasSLAClose(ctx context.Context, client bot.Client, logWriter log.Logger, incident model.Incident) bool {
	if incident.EndTimestamp != nil {
		now := time.Now()
		endTS := incident.EndTimestamp
		diffHours := now.Sub(*endTS)
		if int(diffHours.Hours()) <= config.Env.SLAHoursToClose {
			logWriter.Info(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("SLAHoursToClose"),
				log.NewValue("incident.Status", incident.Status),
				log.NewValue("incident.EndTimestamp", incident.EndTimestamp),
				log.NewValue("SLAHoursToClose", config.Env.SLAHoursToClose),
				log.NewValue("diffHours", diffHours),
			)
			return true
		}
	}

	return false
}
