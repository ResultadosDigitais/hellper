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
	tokenRequest := r.FormValue("token")

	if !authenticateRequest(tokenRequest) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buf.ReadFrom(r.Body)
	body := buf.String()
	logger.Info(
		ctx,
		"handler/pauseNotify.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Info(
		ctx,
		"handler/pauseNotify.ServeHTTP Form",
		formValues...,
	)

	channelID := r.FormValue("channel_id")
	userID := r.FormValue("user_id")
	triggerID := r.FormValue("trigger_id")

	err := commands.PauseNotifyIncidentDialog(ctx, logger, client, repository, channelID, userID, triggerID)
	if err != nil {
		logger.Error(
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
