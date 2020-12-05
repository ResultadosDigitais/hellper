package handler

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
)

func newHandlerClose(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "close", closeIncident, endpoint.NewDefaultSlackErrorHandler())
}

func closeIncident(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	return commands.CloseIncidentDialog(ctx, app, slackParams.ChannelID, slackParams.UserID, slackParams.TriggerID)
}
