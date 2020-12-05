package handler

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
)

func newHandlerStatus(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "status", showIncidentStatus, endpoint.NewDefaultSlackErrorHandler())
}

func showIncidentStatus(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	go func(ctx context.Context) {
		commands.ShowStatus(ctx, app, slackParams.ChannelID, slackParams.UserID)
	}(context.Background())

	return nil
}
