package notify

import (
	"context"
	"errors"
	"hellper/internal"
	"hellper/internal/log"
)

func channelsNotify(ctx context.Context) {

	var (
		to, msg string
	)

	logger.Info(ctx, log.Trace(), log.Action("running"))

	if arg.statusFlag == "" {
		logger.Error(ctx, log.Trace(), log.NewValue("error", errors.New("Must have a status")))
		return
	}

	if arg.toFlag != "" {
		logger.Error(ctx, log.Trace(), log.NewValue("error", errors.New("Forbidden to use the --to option with --type=channels")))
		return
	}

	repository := internal.NewRepository(logger)

	incidents, err := repository.ListActiveIncidents(ctx)
	if err != nil {
		logger.Error(ctx, log.Trace(), log.NewValue("error", err))
	}

	for _, incident := range incidents {

		if arg.toFlag != "" {
			to = arg.toFlag
		} else {
			to = incident.ChannelId
		}

		if arg.msgFlag != "" {
			msg = arg.msgFlag
		} else {
			msg = statusNotify(incident)
		}

		if arg.statusFlag == incident.Status || arg.statusFlag == "all" {
			logger.Info(ctx, log.Trace(), log.Action("notify_job"), log.NewValue("incident", incident))
			err = send(to, msg)
			if err != nil {
				logger.Error(ctx, log.Trace(), log.NewValue("error", err))
			}
		}
	}

}
