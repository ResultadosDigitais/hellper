package notify

import (
	"context"
	"errors"
	"hellper/internal/log"
	"hellper/internal/model"
	"hellper/internal/reminder"
)

func channelsNotify(ctx context.Context) {

	var msg string

	logger.Info(ctx, log.Trace(), log.Action("running"))

	if arg.statusFlag == "" {
		logger.Error(ctx, log.Trace(), log.NewValue("error", errors.New("Must have a status")))
		return
	}

	if arg.toFlag != "" {
		logger.Error(ctx, log.Trace(), log.NewValue("error", errors.New("Forbidden to use the --to option with --type=channels")))
		return
	}

	incidents, err := repository.ListActiveIncidents(ctx)
	if err != nil {
		logger.Error(ctx, log.Trace(), log.NewValue("error", err))
	}

	for _, incident := range incidents {

		if arg.msgFlag != "" {
			msg = arg.msgFlag
		} else {
			msg = statusNotify(incident)
		}

		if arg.statusFlag == incident.Status || arg.statusFlag == "all" {
			notifyChannels(ctx, incident, msg)
		}
	}

}

func notifyChannels(ctx context.Context, incident model.Incident, msg string) {
	if reminder.CanSendNotify(ctx, client, logger, repository, incident) {
		logger.Info(ctx, log.Trace(), log.Action("notify_job"), log.NewValue("incident", incident))
		err := send(incident.ChannelId, msg)
		if err != nil {
			logger.Error(ctx, log.Trace(), log.NewValue("error", err))
		}
	}
}
