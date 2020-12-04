package commands

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"hellper/internal/app"
	"hellper/internal/log"
)

type eventInvoker struct {
	app *app.App
}

var (
	errInvalidCommandLine = errors.New("err_invalid_command_line")
	cmdLineParser         = regexp.MustCompile(`^<@([a-zA-Z0-9]+)>\s*(([a-zA-Z0-9\-_]+)\s*(.*))?$`)
)

func parseCommandLine(cmdLine string) (string, string, error) {
	cmdLine = strings.TrimSpace(cmdLine)
	matches := cmdLineParser.FindStringSubmatch(cmdLine)
	matchesLen := len(matches)
	if matchesLen < 2 {
		return "", "", errInvalidCommandLine
	}
	if matchesLen == 2 {
		return "", "", nil
	}
	if matchesLen == 4 {
		return matches[3], "", nil
	}
	return matches[3], matches[4], nil
}

func newEventInvoker(app *app.App) *eventInvoker {
	return &eventInvoker{
		app: app,
	}
}

func (e *eventInvoker) eventInvoker(ctx context.Context, cmdLine string, event TriggerEvent) error {
	logWriter := e.app.Logger.With(
		log.NewValue("channelID", event.Channel),
		log.NewValue("event", event),
	)

	cmd, args, err := parseCommandLine(cmdLine)
	if err != nil {
		return err
	}

	logWriter.Debug(
		ctx,
		"command/event_invoker.eventInvoker",
		log.NewValue("command", cmd),
		log.NewValue("args", args),
	)

	switch cmd {
	case "", "help":
		help(ctx, e.app, event.Channel)
	case "ping":
		ping(ctx, e.app, event.Channel)
	case "list":
		ListOpenIncidents(ctx, e.app, event)
	}
	return err
}
