package commands

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/log"
)

// EventExecutor represents the executor struct that is used to execute an event command
type EventExecutor struct {
	app *app.App
}

// NewEventExecutor initialize a new EventExecutor type
func NewEventExecutor(app *app.App) *EventExecutor {
	return &EventExecutor{
		app: app,
	}
}

// ExecuteEventCommand calls the invoker passing the command line
func (e *EventExecutor) ExecuteEventCommand(
	ctx context.Context, cmdLine string, event TriggerEvent,
) error {
	e.app.Logger.Info(
		ctx,
		"command/eventexecutor.ExecuteEventCommand",
		log.NewValue("command_line", cmdLine),
		log.NewValue("event", event),
	)

	invoker := newEventInvoker(e.app)
	err := invoker.eventInvoker(ctx, cmdLine, event)
	if err != nil {
		e.app.Logger.Error(
			ctx,
			"command/eventexecutor.ExecuteEventCommand command_result",
			log.NewValue("command_line", cmdLine),
			log.NewValue("error", err),
		)
		return err
	}
	return nil
}
