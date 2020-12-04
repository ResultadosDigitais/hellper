package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/log"
)

type handlerPauseNotify struct {
	app *app.App
}

func newHandlerPauseNotify(app *app.App) *handlerPauseNotify {
	return &handlerPauseNotify{
		app: app,
	}
}

func (h *handlerPauseNotify) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()

		buf        bytes.Buffer
		formValues []log.Value
	)

	r.ParseForm()
	buf.ReadFrom(r.Body)
	body := buf.String()
	h.app.Logger.Debug(
		ctx,
		"handler/pauseNotify.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	h.app.Logger.Debug(
		ctx,
		"handler/pauseNotify.ServeHTTP Form",
		formValues...,
	)

	channelID := r.FormValue("channel_id")
	userID := r.FormValue("user_id")
	triggerID := r.FormValue("trigger_id")

	err := commands.PauseNotifyIncidentDialog(ctx, h.app, channelID, userID, triggerID)
	if err != nil {
		h.app.Logger.Error(
			ctx,
			"handler/pauseNotify.ServeHTTP PauseNotifyIncidentDialog error",
			log.NewValue("channelID", channelID),
			log.NewValue("triggerID", triggerID),
			log.NewValue("error", err),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
