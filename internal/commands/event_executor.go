package commands

import (
	"context"

	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
)

// EventExecutor represents the executor struct that is used to execute an event command
type EventExecutor struct {
	logger     log.Logger
	client     bot.Client
	repository model.Repository
}

// NewEventExecutor initialize a new EventExecutor type
func NewEventExecutor(logger log.Logger, client bot.Client, repository model.Repository) *EventExecutor {
	return &EventExecutor{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

// ExecuteEventCommand calls the invoker passing the command line
func (e *EventExecutor) ExecuteEventCommand(
	ctx context.Context, cmdLine string, event TriggerEvent,
) error {
	e.logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("command_line", cmdLine),
		log.NewValue("event", event),
	)

	invoker := newEventInvoker(e.logger, e.client, e.repository)
	err := invoker.eventInvoker(ctx, cmdLine, event)
	if err != nil {
		e.logger.Error(
			ctx,
			log.Trace(),
			log.Action("invoker.eventInvoker"),
			log.Reason(err.Error()),
			log.NewValue("command_line", cmdLine),
		)
		return err
	}
	return nil
}
