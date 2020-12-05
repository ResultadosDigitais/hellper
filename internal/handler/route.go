package handler

import (
	"fmt"
	"net/http"
	"path"

	"hellper/internal/app"
	"hellper/internal/bot"
)

var (
	openHandler        http.Handler
	editHandler        http.Handler
	eventsHandler      http.Handler
	interactiveHandler http.Handler
	statusHandler      http.Handler
	addStatusHandler   http.Handler
	closeHandler       http.Handler
	cancelHandler      http.Handler
	resolveHandler     http.Handler
	datesHandler       http.Handler
	pauseNotifyHandler http.Handler
)

func init() {
	dependencies := app.NewApp()
	openHandler = newHandlerOpen(&dependencies)
	editHandler = newHandlerEdit(&dependencies)
	eventsHandler = newHandlerEvents(&dependencies)
	interactiveHandler = newHandlerInteractive(&dependencies)
	statusHandler = newHandlerStatus(&dependencies)
	addStatusHandler = newHandlerAddStatus(&dependencies)
	datesHandler = newHandlerDates(&dependencies)
	closeHandler = newHandlerClose(&dependencies)
	cancelHandler = newHandlerCancel(&dependencies)
	resolveHandler = newHandlerResolve(&dependencies)
	pauseNotifyHandler = newHandlerPauseNotify(&dependencies)

	// commands.StartAllReminderJobs(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository)
}

// NewHandlerRoute handles the http requests received and calls the correct handler.
func NewHandlerRoute() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		lastPath := path.Base(r.URL.Path)

		switch lastPath {
		case "healthz":
			fmt.Fprintf(w, "I'm working!!")
		case "events":
			bot.VerifyRequests(r, w, eventsHandler)
		case "open":
			bot.VerifyRequests(r, w, openHandler)
		case "edit":
			bot.VerifyRequests(r, w, editHandler)
		case "interactive":
			bot.VerifyRequests(r, w, interactiveHandler)
		case "status":
			bot.VerifyRequests(r, w, statusHandler)
		case "add-status":
			bot.VerifyRequests(r, w, addStatusHandler)
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
