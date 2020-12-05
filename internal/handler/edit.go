package handler

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
)

func newHandlerEdit(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "edit", editIncident, endpoint.NewDefaultSlackErrorHandler())
}

func editIncident(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	return commands.OpenEditIncidentDialog(ctx, app, slackParams.ChannelID, slackParams.TriggerID)
}
