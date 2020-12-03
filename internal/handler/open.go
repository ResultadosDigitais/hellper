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
	logger            log.Logger
	client            bot.Client
	repository        model.IncidentRepository
	serviceRepository model.ServiceRepository
}

func newHandlerOpen(
	logger log.Logger,
	client bot.Client,
	repository model.IncidentRepository,
	serviceRepository model.ServiceRepository,
) *handlerOpen {
	return &handlerOpen{
		logger:            logger,
		client:            client,
		repository:        repository,
		serviceRepository: serviceRepository,
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

	err := commands.OpenStartIncidentDialog(ctx, h.client, h.serviceRepository, triggerID)
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
