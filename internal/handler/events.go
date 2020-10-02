package handler

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"

	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"

	"github.com/slack-go/slack/slackevents"
)

var msgsCache = map[string]struct{}{}

type handlerEvents struct {
	logger     log.Logger
	client     bot.Client
	repository model.Repository
}

func stringSha1(v string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(v)))
}

func newHandlerEvents(logger log.Logger, client bot.Client, repository model.Repository) *handlerEvents {
	return &handlerEvents{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (h *handlerEvents) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		buf    bytes.Buffer
		logger = h.logger
	)
	buf.ReadFrom(r.Body)
	body := buf.String()
	logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("requestbody", body),
	)

	event, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionNoVerifyToken(),
	)
	if err != nil {
		logger.Error(
			ctx,
			log.Trace(),
			log.Action("slackevents.ParseEvent"),
			log.Reason(err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case slackevents.CallbackEvent:
		logger.Info(
			ctx,
			log.Trace(),
			log.NewValue("event", event),
		)

		// Temporary memory cache
		if _, exists := msgsCache[stringSha1(body)]; exists {
			logger.Info(
				ctx,
				log.Trace(),
				log.NewValue("event", event),
				log.NewValue("message", msgsCache),
			)
			return
		}
		msgsCache[stringSha1(body)] = struct{}{}
		logger.Info(
			ctx,
			log.Trace(),
			log.NewValue("event", event),
			log.NewValue("message", msgsCache),
		)

		err = replyCallbackEvent(ctx, h.logger, h.client, h.repository, event)
		if err != nil {
			logger.Error(
				ctx,
				log.Trace(),
				log.Action("replyCallbackEvent"),
				log.Reason(err.Error()),
				log.NewValue("event", event),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	case slackevents.URLVerification:
		logger.Info(
			ctx,
			log.Trace(),
			log.NewValue("event", event),
		)

		var resp slackevents.ChallengeResponse
		err = json.NewDecoder(&buf).Decode(&resp)
		if err != nil {
			logger.Error(
				ctx,
				log.Trace(),
				log.Action("json.NewDecoder"),
				log.Reason(err.Error()),
				log.NewValue("event", event),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "text")
		fmt.Fprintf(w, "%s", resp.Challenge)

		logger.Info(
			ctx,
			log.Trace(),
			log.NewValue("event", event),
			log.NewValue("challenge", resp.Challenge),
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
