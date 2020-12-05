package handler

import (
	"context"
	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
)

func newHandlerOpen(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "open", openIncident, endpoint.NewDefaultSlackErrorHandler())
}

func openIncident(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	return commands.OpenStartIncidentDialog(ctx, app, slackParams.UserID, slackParams.TriggerID)
}
