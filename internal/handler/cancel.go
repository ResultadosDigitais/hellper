package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/log"
)

type handlerCancel struct {
	app *app.App
}

func newHandlerCancel(app *app.App) *handlerCancel {
	return &handlerCancel{
		app: app,
	}
}

func (h *handlerCancel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		"handler/cancel.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Debug(
		ctx,
		"handler/cancel.ServeHTTP Form",
		formValues...,
	)

	tiggerID := r.FormValue("trigger_id")
	channelID := r.FormValue("channel_id")
	userID := r.FormValue("user_id")

	err := commands.OpenCancelIncidentDialog(ctx, h.app, channelID, userID, tiggerID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("commands.OpenCancelIncidentDialog"),
			log.NewValue("error", err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
