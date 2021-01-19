package reminder

import (
	"context"
	"hellper/internal/bot"
	"hellper/internal/config"
	"hellper/internal/log"
	"hellper/internal/model"
	"time"
)

func hasSLAClose(ctx context.Context, client bot.Client, logger log.Logger, incident model.Incident) bool {
	if incident.EndedAt != nil {
		now := time.Now()
		endTS := incident.EndedAt
		diffHours := now.Sub(*endTS)
		if int(diffHours.Hours()) <= config.Env.SLAHoursToClose {
			logger.Info(
				ctx,
				log.Trace(),
				log.Action("do_not_notify"),
				log.Reason("SLAHoursToClose"),
				log.NewValue("channelID", incident.ChannelID),
				log.NewValue("channelName", incident.ChannelName),
				log.NewValue("incident.Status", incident.Status),
				log.NewValue("incident.EndedAt", incident.EndedAt),
				log.NewValue("SLAHoursToClose", config.Env.SLAHoursToClose),
				log.NewValue("diffHours", diffHours),
			)
			return true
		}
	}

	return false
}
