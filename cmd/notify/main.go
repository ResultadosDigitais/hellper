package main

import (
	"context"
	"hellper/internal"
	"hellper/internal/log"
)

var (
	logger     = internal.NewLogger()
	client     = internal.NewClient(logger)
	repository = internal.NewRepository(logger)
)

func main() {
	ctx := context.Background()
	logger.Info(ctx, log.Trace(), log.Action("running"))

	incidents, err := repository.ListActiveIncidents(ctx)
	if err != nil {
		logger.Error(ctx, log.Trace(), log.NewValue("error", err))
	}

	for _, incident := range incidents {
		logger.Info(ctx, log.Trace(), log.Action("ListActiveIncidents"), log.NewValue("incident", incident))
	}
}
