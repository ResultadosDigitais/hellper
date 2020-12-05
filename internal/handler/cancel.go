package handler

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
)

func newHandlerCancel(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "cancel", cancelIncident, endpoint.NewDefaultSlackErrorHandler())
}

func cancelIncident(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	return commands.OpenCancelIncidentDialog(ctx, app, slackParams.ChannelID, slackParams.UserID, slackParams.TriggerID)
}
