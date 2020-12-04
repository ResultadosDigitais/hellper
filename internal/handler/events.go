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
	h.app.Logger.Info(
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

	switch event.Type {
	case slackevents.CallbackEvent:
		h.app.Logger.Info(
			ctx,
			"handler/events.ParseEvent CallbackEvent",
			log.NewValue("event", event),
		)

		// Temporary memory cache
		if _, exists := msgsCache[stringSha1(body)]; exists {
			h.app.Logger.Info(
				ctx,
				"handler/events.ParseEvent duplicated_message",
				log.NewValue("event", event),
				log.NewValue("message", msgsCache),
			)
			return
		}
		msgsCache[stringSha1(body)] = struct{}{}
		h.app.Logger.Info(
			ctx,
			"handler/events.ParseEvent deduplication_message_added",
			log.NewValue("event", event),
			log.NewValue("message", msgsCache),
		)

		err = replyCallbackEvent(ctx, h.app, event)
		if err != nil {
			h.app.Logger.Error(
				ctx,
				"handler/events.ParseEvent callback_event_error",
				log.NewValue("event", event),
				log.NewValue("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	case slackevents.URLVerification:
		h.app.Logger.Info(
			ctx,
			"handler/events.ParseEvent URLVerification",
			log.NewValue("event", event),
		)

		var resp slackevents.ChallengeResponse
		err = json.NewDecoder(&buf).Decode(&resp)
		if err != nil {
			h.app.Logger.Error(
				ctx,
				"handler/events.ParseEvent Decode error",
				log.NewValue("event", event),
				log.NewValue("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "text")
		fmt.Fprintf(w, "%s", resp.Challenge)

		h.app.Logger.Info(
			ctx,
			"handler/events.ParseEvent challenge",
			log.NewValue("event", event),
			log.NewValue("challenge", resp.Challenge),
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
