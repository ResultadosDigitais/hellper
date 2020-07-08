package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"hellper/internal/bot"
	calendar "hellper/internal/calendar"
	"hellper/internal/commands"
	filestorage "hellper/internal/file_storage"
	"hellper/internal/log"
	"hellper/internal/model"
)

type handlerInteractive struct {
	logger      log.Logger
	client      bot.Client
	repository  model.Repository
	fileStorage filestorage.Driver
	calendar    calendar.Calendar
}

func newHandlerInteractive(logger log.Logger, client bot.Client, repository model.Repository, fileStorage filestorage.Driver, calendar calendar.Calendar) *handlerInteractive {
	return &handlerInteractive{
		logger:      logger,
		client:      client,
		repository:  repository,
		fileStorage: fileStorage,
		calendar:    calendar,
	}
}

func (h *handlerInteractive) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		"handler/interactive.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}
	logger.Info(
		ctx,
		"handler/interactive.ServeHTTP Form",
		formValues...,
	)

	formPayload := r.FormValue("payload")

	dialogSubmission := bot.DialogSubmission{}
	json.Unmarshal([]byte(formPayload), &dialogSubmission)

	logger.Info(
		ctx,
		"handler/interactive.ServeHTTP dialogSubmission",
		log.NewValue("dialogSubmission", dialogSubmission),
	)

	tokenRequest := dialogSubmission.Token
	if !authenticateRequest(tokenRequest) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	callbackID := dialogSubmission.CallbackID
	var err error

	switch callbackID {
	case "inc-close":
		err = commands.CloseIncidentByDialog(ctx, h.client, h.logger, h.repository, dialogSubmission)
	case "inc-cancel":
		err = commands.CancelIncidentByDialog(ctx, h.client, h.logger, h.repository, dialogSubmission)
	case "inc-open":
		err = commands.StartIncidentByDialog(ctx, h.client, h.logger, h.repository, h.fileStorage, dialogSubmission)
	case "inc-resolve":
		err = commands.ResolveIncidentByDialog(ctx, h.client, h.logger, h.repository, h.calendar, dialogSubmission)
	case "inc-dates":
		err = commands.UpdateDatesByDialog(ctx, h.client, h.logger, h.repository, dialogSubmission)
	case "inc-pausenotify":
		err = commands.PauseNotifyIncidentByDialog(ctx, h.client, h.logger, h.repository, dialogSubmission)
	default:
		commands.PostErrorAttachment(
			ctx,
			h.client,
			h.logger,
			dialogSubmission.Channel.ID,
			dialogSubmission.User.ID,
			"invalid command, "+callbackID,
		)
		logger.Error(
			ctx,
			"handler/interactive.ServeHTTP invalid_callbackID",
			log.NewValue("dialogSubmission", dialogSubmission),
		)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		logger.Error(
			ctx,
			"handler/interactive.ServeHTTP proccess_submit_dialog_error",
			log.NewValue("error", err),
		)

		commands.PostErrorAttachment(ctx, h.client, h.logger, dialogSubmission.Channel.ID, dialogSubmission.User.ID, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
