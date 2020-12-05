package handler

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
)

func newHandlerResolve(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "resolve", resolveIncident, endpoint.NewDefaultSlackErrorHandler())
}

func resolveIncident(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	return commands.ResolveIncidentDialog(app, slackParams.TriggerID)
}
