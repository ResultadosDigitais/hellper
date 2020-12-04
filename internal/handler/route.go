package handler

import (
	"fmt"
	"net/http"
	"path"

	"hellper/internal"
	"hellper/internal/bot"
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
	dependencies := internal.NewApp()
	openHandler = newHandlerOpen(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository, dependencies.ServiceRepository)
	eventsHandler = newHandlerEvents(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository)
	interactiveHandler = newHandlerInteractive(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository,
		dependencies.FileStorage, dependencies.Calendar)
	statusHandler = newHandlerStatus(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository)
	datesHandler = newHandlerDates(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository)
	closeHandler = newHandlerClose(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository)
	cancelHandler = newHandlerCancel(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository)
	resolveHandler = newHandlerResolve(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository)
	pauseNotifyHandler = newHandlerPauseNotify(dependencies.Logger, dependencies.Client, dependencies.IncidentRepository)
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
