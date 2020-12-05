package handler

import (
	"context"

	"hellper/internal/app"
	"hellper/internal/commands"
	"hellper/internal/handler/endpoint"
)

func newHandlerPauseNotify(app *app.App) *endpoint.Endpoint {
	return endpoint.NewSlackEndpoint(app, "pauseNotifications", pauseIncidentNotifications, endpoint.NewDefaultSlackErrorHandler())
}

func pauseIncidentNotifications(ctx context.Context, app *app.App, slackParams endpoint.SlackParams, endpointContext *endpoint.Context) error {
	return commands.PauseNotifyIncidentDialog(ctx, app, slackParams.ChannelID, slackParams.UserID, slackParams.TriggerID)
}
