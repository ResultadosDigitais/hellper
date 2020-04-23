package handler

import (
	"fmt"
	"net/http"
	"path"

	"hellper/internal"
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
)

func init() {
	logger, client, repository, fileStorage := internal.New()
	openHandler = newHandlerOpen(logger, client, repository)
	eventsHandler = newHandlerEvents(logger, client, repository)
	interactiveHandler = newHandlerInteractive(logger, client, repository, fileStorage)
	statusHandler = newHandlerStatus(logger, client, repository)
	datesHandler = newHandlerDates(logger, client, repository)
	closeHandler = newHandlerClose(logger, client, repository)
	cancelHandler = newHandlerCancel(logger, client, repository)
	resolveHandler = newHandlerResolve(logger, client, repository)
	commands.StartAllReminderJobs(logger, client, repository)
}

func authenticateRequest(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	requestToken := r.FormValue("token")

	if requestToken != config.Env.VerificationToken {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	w.WriteHeader(http.StatusAccepted)
	return true
}

// NewHandlerRoute handles the http requests received and calls the correct handler.
func NewHandlerRoute() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lastPath := path.Base(r.URL.Path)

		switch lastPath {
		case "healthz":
			fmt.Fprintf(w, "I'm working!!")
		case "envtest":
			fmt.Fprintf(w, "%+v\n", config.Env.Messages)
		case "events":
			eventsHandler.ServeHTTP(w, r)
		case "open":
			openHandler.ServeHTTP(w, r)
			// if authenticateRequest(w, r) {
			// 	openHandler.ServeHTTP(w, r)
			// }
		case "interactive":
			if authenticateRequest(w, r) {
				interactiveHandler.ServeHTTP(w, r)
			}
		case "status":
			if authenticateRequest(w, r) {
				statusHandler.ServeHTTP(w, r)
			}
		case "dates":
			if authenticateRequest(w, r) {
				datesHandler.ServeHTTP(w, r)
			}
		case "close":
			if authenticateRequest(w, r) {
				closeHandler.ServeHTTP(w, r)
			}
		case "cancel":
			w.WriteHeader(http.StatusNotImplemented)
		case "resolve":
			resolveHandler.ServeHTTP(w, r)
		default:
			fmt.Fprintf(w, "invalid path, %s!", lastPath)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
