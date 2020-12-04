package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/log"
)

type handlerClose struct {
	app *app.App
}

func newHandlerClose(app *app.App) *handlerClose {
	return &handlerClose{
		app: app,
	}
}

func (h *handlerClose) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = h.app.Logger

		formValues []log.Value
		buf        bytes.Buffer
	)

	r.ParseForm()
	buf.ReadFrom(r.Body)
	body := buf.String()
	logger.Debug(
		ctx,
		"handler/close.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Debug(
		ctx,
		"handler/close.ServeHTTP Form",
		formValues...,
	)

	triggerID := r.FormValue("trigger_id")
	channelID := r.FormValue("channel_id")
	userID := r.FormValue("user_id")

	err := commands.CloseIncidentDialog(ctx, h.app, channelID, userID, triggerID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("commands.CloseIncidentDialog"),
			log.NewValue("error", err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
