package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/log"
)

type handlerOpen struct {
	app *app.App
}

func newHandlerOpen(
	app *app.App,
) *handlerOpen {
	return &handlerOpen{
		app: app,
	}
}

func (h *handlerOpen) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		"handler/open.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Debug(
		ctx,
		"handler/open.ServeHTTP Form",
		formValues...,
	)

	triggerID := r.FormValue("trigger_id")

	err := commands.OpenStartIncidentDialog(ctx, h.app, triggerID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Reason("OpenStartIncidentDialog"),
			log.NewValue("triggerID", triggerID),
			log.NewValue("error", err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
