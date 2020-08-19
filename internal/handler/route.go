package handler

import (
	"fmt"
	"net/http"
	"path"

	"hellper/internal"
	"hellper/internal/bot"
	"hellper/internal/commands"
	"hellper/internal/config"
)

var (
	openHandler        http.Handler
	eventsHandler      http.Handler
	interactiveHandler http.Handler
	statusHandler      http.Handler
	closeHandler       http.Handler
	cancelHandler      http.Handler
	resolveHandler     http.Handler
	datesHandler       http.Handler
	pauseNotifyHandler http.Handler
)

func init() {
	logger, client, repository, fileStorage, calendar := internal.New()
	openHandler = newHandlerOpen(logger, client, repository)
	eventsHandler = newHandlerEvents(logger, client, repository)
	interactiveHandler = newHandlerInteractive(logger, client, repository, fileStorage, calendar)
	statusHandler = newHandlerStatus(logger, client, repository)
	datesHandler = newHandlerDates(logger, client, repository)
	closeHandler = newHandlerClose(logger, client, repository)
	cancelHandler = newHandlerCancel(logger, client, repository)
	resolveHandler = newHandlerResolve(logger, client, repository)
	pauseNotifyHandler = newHandlerPauseNotify(logger, client, repository)
	commands.StartAllReminderJobs(logger, client, repository)
}

func authenticateRequest(token string) bool {
	if token != config.Env.VerificationToken {
		return false
	}

	return true
}

// NewHandlerRoute handles the http requests received and calls the correct handler.
func NewHandlerRoute() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		lastPath := path.Base(r.URL.Path)

		switch lastPath {
		case "healthz":
			fmt.Fprintf(w, "I'm working!!")
		case "envtest":
			fmt.Fprintf(w, "%+v\n", config.Env.Messages)
		case "events":
			bot.VerifyRequests(r, w, eventsHandler)
		case "open":
			bot.VerifyRequests(r, w, openHandler)
		case "interactive":
			bot.VerifyRequests(r, w, interactiveHandler)
		case "status":
			bot.VerifyRequests(r, w, statusHandler)
		case "dates":
			bot.VerifyRequests(r, w, datesHandler)
		case "close":
			bot.VerifyRequests(r, w, closeHandler)
		case "cancel":
			bot.VerifyRequests(r, w, cancelHandler)
		case "resolve":
			bot.VerifyRequests(r, w, resolveHandler)
		case "pause-notify":
			bot.VerifyRequests(r, w, pauseNotifyHandler)
		default:
			fmt.Fprintf(w, "invalid path, %s!", lastPath)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
