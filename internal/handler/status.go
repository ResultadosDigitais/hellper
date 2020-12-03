package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
)

type handlerStatus struct {
	logger     log.Logger
	client     bot.Client
	repository model.Repository
}

func newHandlerStatus(logger log.Logger, client bot.Client, repository model.Repository) *handlerStatus {
	return &handlerStatus{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (h *handlerStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = h.logger

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

	err := commands.ShowStatus(ctx, h.client, logger, h.repository, channelID, userID)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("commands.ShowStatus"),
			log.Reason(err.Error()),
			log.NewValue("channelID", channelID),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
