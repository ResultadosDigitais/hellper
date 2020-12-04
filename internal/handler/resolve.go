package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/log"
)

type handlerResolve struct {
	app *app.App
}

func newHandlerResolve(app *app.App) *handlerResolve {
	return &handlerResolve{
		app: app,
	}
}

func (h *handlerResolve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		"handler/resolve.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	h.app.Logger.Debug(
		ctx,
		"handler/resolve.ServeHTTP Form",
		formValues...,
	)

	triggerID := r.FormValue("trigger_id")

	err := commands.ResolveIncidentDialog(h.app, triggerID)
	if err != nil {
		h.app.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason("ResolveIncidentDialog"),
			log.NewValue("triggerID", triggerID),
			log.NewValue("error", err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
