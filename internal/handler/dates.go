package handler

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
)

func newHandlerDates(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "updateDates", updateIncidentDates, endpoint.NewDefaultSlackErrorHandler())
}

func updateIncidentDates(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	return commands.CloseIncidentDialog(ctx, app, slackParams.ChannelID, slackParams.UserID, slackParams.TriggerID)
}
