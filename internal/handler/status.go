package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/log"
)

type handlerStatus struct {
	app *app.App
}

func newHandlerStatus(app *app.App) *handlerStatus {
	return &handlerStatus{
		app: app,
	}
}

func (h *handlerStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		"handler/status.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	h.app.Logger.Debug(
		ctx,
		"handler/status.ServeHTTP Form",
		formValues...,
	)

	channelID := r.FormValue("channel_id")
	userID := r.FormValue("user_id")

	err := commands.ShowStatus(ctx, h.app, channelID, userID)
	if err != nil {
		h.app.Logger.Error(
			ctx,
			"handler/status.ServeHTTP ShowStatus error",
			log.NewValue("channelID", channelID),
			log.NewValue("error", err),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
