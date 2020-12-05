package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
	"hellper/internal/log"
)

func newHandlerInteractive(app *app.App) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		App:          app,
		Name:         "interactive",
		ErrorHandler: endpoint.NewDefaultSlackErrorHandler(),
		Function:     processInteractiveEvent,
	}
}

func processInteractiveEvent(ctx context.Context, app *app.App, endpointContext *endpoint.Context) error {
	var payload string

	endpointContext.ReadForm().Read(&payload, "payload")
	dialogSubmission := bot.DialogSubmission{}
	json.Unmarshal([]byte(payload), &dialogSubmission)

	app.Logger.Debug(
		ctx,
		"handler/interactive.ServeHTTP dialogSubmission",
		log.NewValue("dialogSubmission", dialogSubmission),
	)

	callbackID := dialogSubmission.CallbackID
	err := processEvent(ctx, app, endpointContext, callbackID, dialogSubmission)

	if err != nil {
		app.Logger.Error(
			ctx,
			"handler/interactive.ServeHTTP proccess_submit_dialog_error",
			log.NewValue("error", err),
		)

		commands.PostErrorAttachment(ctx, app, dialogSubmission.Channel.ID, dialogSubmission.User.ID, err.Error())
	}

	endpointContext.Status(http.StatusNoContent)
	return nil
}

func processEvent(
	ctx context.Context, app *app.App, endpointContext *endpoint.Context,
	callbackID string, dialogSubmission bot.DialogSubmission,
) error {
	switch callbackID {
	case "inc-close":
		return commands.CloseIncidentByDialog(ctx, app, dialogSubmission)
	case "inc-cancel":
		return commands.CancelIncidentByDialog(ctx, app, dialogSubmission)
	case "inc-open":
		return commands.StartIncidentByDialog(ctx, app, dialogSubmission)
	case "inc-edit":
		return commands.EditIncidentByDialog(ctx, app, dialogSubmission)
	case "inc-resolve":
		return commands.ResolveIncidentByDialog(ctx, app, dialogSubmission)
	case "inc-dates":
		return commands.UpdateDatesByDialog(ctx, app, dialogSubmission)
	case "inc-pausenotify":
		return commands.PauseNotifyIncidentByDialog(ctx, app, dialogSubmission)
	default:
		commands.PostErrorAttachment(
			ctx,
			app,
			dialogSubmission.Channel.ID,
			dialogSubmission.User.ID,
			"invalid command, "+callbackID,
		)

		app.Logger.Error(
			ctx,
			"handler/interactive.ServeHTTP invalid_callbackID",
			log.NewValue("dialogSubmission", dialogSubmission),
		)

		endpointContext.Status(http.StatusBadRequest)
	}

	return nil
}
