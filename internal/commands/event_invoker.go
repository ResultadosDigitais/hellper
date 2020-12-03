package commands

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"hellper/internal/bot"
	"hellper/internal/log"
	"hellper/internal/model"
)

type eventInvoker struct {
	logger     log.Logger
	client     bot.Client
	repository model.Repository
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

func newEventInvoker(
	logger log.Logger, client bot.Client, repository model.Repository,
) *eventInvoker {
	return &eventInvoker{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (e *eventInvoker) eventInvoker(ctx context.Context, cmdLine string, event TriggerEvent) error {
	var (
		logger     = e.logger
		client     = e.client
		repository = e.repository
	)

	cmd, args, err := parseCommandLine(cmdLine)
	if err != nil {
		return err
	}

	logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("command", cmd),
		log.NewValue("args", args),
		log.NewValue("event", event),
	)

	switch cmd {
	case "", "help":
		help(ctx, client, logger, event.Channel)
	case "ping":
		ping(ctx, client, logger, event.Channel)
	case "list":
		ListOpenIncidents(ctx, client, logger, repository, event)
	}
	return err
}
