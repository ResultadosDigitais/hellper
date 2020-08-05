package notify

import (
	"context"
	"errors"
	"hellper/internal"
	"hellper/internal/log"
	"strings"
)

func reportNotify(ctx context.Context) {

	logger.Info(ctx, log.Trace(), log.Action("running"))

	if arg.statusFlag == "" {
		logger.Error(ctx, log.Trace(), log.NewValue("error", errors.New("Must have a status")))
		return
	}

	if arg.toFlag == "" {
		logger.Error(ctx, log.Trace(), log.NewValue("error", errors.New("Must have a destination")))
		return
	}

	if arg.msgFlag != "" {
		logger.Error(ctx, log.Trace(), log.NewValue("error", errors.New("Forbidden to use the --msg option with --type=report")))
		return
	}

	repository := internal.NewRepository(logger)

	incidents, err := repository.ListActiveIncidents(ctx)
	if err != nil {
		logger.Error(ctx, log.Trace(), log.NewValue("error", err))
	}

	var notify strings.Builder
	notify.WriteString(":mega: *Incident Reporting:*\n")

	for _, incident := range incidents {
		if arg.statusFlag == incident.Status || arg.statusFlag == "all" {
			logger.Info(ctx, log.Trace(), log.Action("notify_job"), log.NewValue("incident", incident))
			notify.WriteString("*<#" + incident.ChannelId + ">* - ")
			notify.WriteString("Status: `" + incident.Status + "` - ")
			notify.WriteString("Commander: <@" + incident.CommanderId + ">\n")
		}
	}

	err = send(arg.toFlag, notify.String())
	if err != nil {
		logger.Error(ctx, log.Trace(), log.NewValue("error", err))
	}

}
