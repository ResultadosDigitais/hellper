package handler

import (
	"context"
	"fmt"
	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
	"strings"
)

func newHandlerAddStatus(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "addStatus", addIncidentStatus, endpoint.NewDefaultSlackErrorHandler())
}

func addIncidentStatus(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	var (
		message  string
		userName string
	)

	endpointContext.ReadForm().
		Read(&message, "text").
		Read(&userName, "user_name")

	if !isValidMessage(message) {
		return fmt.Errorf("Your message must have at least one character")
	}

	go func(ctx context.Context) {
		commands.AddStatus(ctx, app, slackParams.ChannelID, slackParams.UserID, userName, message)
	}(context.Background())

	return nil
}

func isValidMessage(message string) bool {
	return len(strings.Trim(message, " ")) > 0
}
