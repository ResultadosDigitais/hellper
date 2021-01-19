package reminder

import (
	"context"
	"hellper/internal/log"
	"hellper/internal/model"
	"time"
)

func hasSnoozedUntil(ctx context.Context, logger log.Logger, incident model.Incident) bool {
	snoozedUntil := incident.SnoozedUntil
	if snoozedUntil.Time.Unix() > time.Now().Unix() {
		logger.Info(
			ctx,
			log.Trace(),
			log.Action("do_not_notify"),
			log.Reason("isPaused"),
			log.NewValue("channelID", incident.ChannelID),
			log.NewValue("channelName", incident.ChannelName),
			log.NewValue("snoozedUntil", snoozedUntil.Time),
		)
		return true
	}

	return false
}
