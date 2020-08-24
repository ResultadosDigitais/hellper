package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
)

type handlerResolve struct {
	logger     log.Logger
	client     bot.Client
	repository model.Repository
}

func newHandlerResolve(logger log.Logger, client bot.Client, repository model.Repository) *handlerResolve {
	return &handlerResolve{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (h *handlerResolve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = h.logger

		formValues []log.Value
		buf        bytes.Buffer
	)

	r.ParseForm()
	buf.ReadFrom(r.Body)
	body := buf.String()
	logger.Info(
		ctx,
		"handler/resolve.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Info(
		ctx,
		"handler/resolve.ServeHTTP Form",
		formValues...,
	)

	triggerID := r.FormValue("trigger_id")

	err := commands.ResolveIncidentDialog(h.client, triggerID)
	if err != nil {
		logger.Error(
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
