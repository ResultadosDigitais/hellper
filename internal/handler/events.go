package handler

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"

	"hellper/internal/app"
	"hellper/internal/log"

	"github.com/slack-go/slack/slackevents"
)

var msgsCache = map[string]struct{}{}

type handlerEvents struct {
	app *app.App
}

func stringSha1(v string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(v)))
}

func newHandlerEvents(app *app.App) *handlerEvents {
	return &handlerEvents{
		app: app,
	}
}

func (h *handlerEvents) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		buf bytes.Buffer
	)
	buf.ReadFrom(r.Body)
	body := buf.String()
	h.app.Logger.Debug(
		ctx,
		"handler/events.ServeHTTP",
		log.NewValue("requestbody", body),
	)

	event, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionNoVerifyToken(),
	)
	if err != nil {
		h.app.Logger.Error(
			ctx,
			"handler/events.ServeHTTP ParseEvent error",
			log.NewValue("error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logWriter := h.app.Logger.With(
		log.NewValue("event", event),
	)

	switch event.Type {
	case slackevents.CallbackEvent:
		logWriter.Debug(
			ctx,
			"handler/events.ParseEvent CallbackEvent",
		)

		// Temporary memory cache
		if _, exists := msgsCache[stringSha1(body)]; exists {
			logWriter.Debug(
				ctx,
				"handler/events.ParseEvent duplicated_message",
				log.NewValue("message", msgsCache),
			)
			return
		}
		msgsCache[stringSha1(body)] = struct{}{}
		logWriter.Debug(
			ctx,
			"handler/events.ParseEvent deduplication_message_added",
			log.NewValue("message", msgsCache),
		)

		err = replyCallbackEvent(ctx, h.app, event)
		if err != nil {
			logWriter.Error(
				ctx,
				"handler/events.ParseEvent callback_event_error",
				log.NewValue("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	case slackevents.URLVerification:
		logWriter.Debug(
			ctx,
			"handler/events.ParseEvent URLVerification",
			log.NewValue("event", event),
		)

		var resp slackevents.ChallengeResponse
		err = json.NewDecoder(&buf).Decode(&resp)
		if err != nil {
			logWriter.Error(
				ctx,
				"handler/events.ParseEvent Decode error",
				log.NewValue("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "text")
		fmt.Fprintf(w, "%s", resp.Challenge)

		logWriter.Debug(
			ctx,
			"handler/events.ParseEvent challenge",
			log.NewValue("challenge", resp.Challenge),
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
