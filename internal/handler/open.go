package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
)

type handlerOpen struct {
	logger     log.Logger
	client     bot.Client
	repository model.Repository
}

func newHandlerOpen(logger log.Logger, client bot.Client, repository model.Repository) *handlerOpen {
	return &handlerOpen{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (h *handlerOpen) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		"handler/open.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Info(
		ctx,
		"handler/open.ServeHTTP Form",
		formValues...,
	)

	triggerID := r.FormValue("trigger_id")

	err := commands.OpenStartIncidentDialog(h.client, triggerID)
	if err != nil {
		logger.Error(
			ctx,
			"handler/open.ServeHTTP OpenStartIncidentDialog error",
			log.NewValue("error", err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
