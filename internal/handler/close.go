package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
)

type handlerClose struct {
	logger     log.Logger
	client     bot.Client
	repository model.Repository
}

func newHandlerClose(logger log.Logger, client bot.Client, repository model.Repository) *handlerClose {
	return &handlerClose{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (h *handlerClose) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		log.Trace(),
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Info(
		ctx,
		log.Trace(),
		formValues...,
	)

	triggerID := r.FormValue("trigger_id")
	channelID := r.FormValue("channel_id")
	userID := r.FormValue("user_id")

	err := commands.CloseIncidentDialog(ctx, logger, h.client, h.repository, channelID, userID, triggerID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("commands.CloseIncidentDialog"),
			log.Reason(err.Error()),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
