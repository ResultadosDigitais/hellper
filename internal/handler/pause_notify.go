package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
)

type handlerPauseNotify struct {
	logger     log.Logger
	client     bot.Client
	repository model.Repository
}

func newHandlerPauseNotify(logger log.Logger, client bot.Client, repository model.Repository) *handlerPauseNotify {
	return &handlerPauseNotify{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (h *handlerPauseNotify) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		logger     = h.logger
		client     = h.client
		repository = h.repository

		buf        bytes.Buffer
		formValues []log.Value
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

	channelID := r.FormValue("channel_id")
	userID := r.FormValue("user_id")
	triggerID := r.FormValue("trigger_id")

	err := commands.PauseNotifyIncidentDialog(ctx, logger, client, repository, channelID, userID, triggerID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("commands.PauseNotifyIncidentDialog"),
			log.Reason(err.Error()),
			log.NewValue("channelID", channelID),
			log.NewValue("triggerID", triggerID),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
