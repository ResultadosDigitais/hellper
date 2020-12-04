package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
)

type handlerInteractive struct {
	app *app.App
}

func newHandlerInteractive(app *app.App) *handlerInteractive {
	return &handlerInteractive{
		app: app,
	}
}

func (h *handlerInteractive) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()

		formValues []log.Value
		buf        bytes.Buffer
	)

	r.ParseForm()

	buf.ReadFrom(r.Body)
	body := buf.String()
	h.app.Logger.Debug(
		ctx,
		"handler/interactive.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	h.app.Logger.Debug(
		ctx,
		"handler/interactive.ServeHTTP Form",
		formValues...,
	)

	formPayload := r.FormValue("payload")

	dialogSubmission := bot.DialogSubmission{}
	json.Unmarshal([]byte(formPayload), &dialogSubmission)

	h.app.Logger.Debug(
		ctx,
		"handler/interactive.ServeHTTP dialogSubmission",
		log.NewValue("dialogSubmission", dialogSubmission),
	)

	callbackID := dialogSubmission.CallbackID
	var err error

	switch callbackID {
	case "inc-close":
		err = commands.CloseIncidentByDialog(ctx, h.app, dialogSubmission)
	case "inc-cancel":
		err = commands.CancelIncidentByDialog(ctx, h.app, dialogSubmission)
	case "inc-open":
		err = commands.StartIncidentByDialog(ctx, h.app, dialogSubmission)
	case "inc-resolve":
		err = commands.ResolveIncidentByDialog(ctx, h.app, dialogSubmission)
	case "inc-dates":
		err = commands.UpdateDatesByDialog(ctx, h.app, dialogSubmission)
	case "inc-pausenotify":
		err = commands.PauseNotifyIncidentByDialog(ctx, h.app, dialogSubmission)
	default:
		commands.PostErrorAttachment(
			ctx,
			h.app,
			dialogSubmission.Channel.ID,
			dialogSubmission.User.ID,
			"invalid command, "+callbackID,
		)
		h.app.Logger.Error(
			ctx,
			"handler/interactive.ServeHTTP invalid_callbackID",
			log.NewValue("dialogSubmission", dialogSubmission),
		)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		h.app.Logger.Error(
			ctx,
			"handler/interactive.ServeHTTP proccess_submit_dialog_error",
			log.NewValue("error", err),
		)

		commands.PostErrorAttachment(ctx, h.app, dialogSubmission.Channel.ID, dialogSubmission.User.ID, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
