package handler

import (
	"bytes"
	"net/http"

	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/log"
	"hellper/internal/model"
)

type handlerCancel struct {
	logger     log.Logger
	client     bot.Client
	repository model.IncidentRepository
}

func newHandlerCancel(logger log.Logger, client bot.Client, repository model.IncidentRepository) *handlerCancel {
	return &handlerCancel{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (h *handlerCancel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		"handler/cancel.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Info(
		ctx,
		"handler/cancel.ServeHTTP Form",
		formValues...,
	)

	tiggerID := r.FormValue("trigger_id")
	channelID := r.FormValue("channel_id")
	userID := r.FormValue("user_id")

	err := commands.OpenCancelIncidentDialog(ctx, logger, h.client, h.repository, channelID, userID, tiggerID)
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
