package endpoint

import (
	"context"
	"fmt"
	"hellper/internal/app"
	"hellper/internal/log"
	"net/http"
)

// SlackEndpointFunction is the function that will process the slack request
type SlackEndpointFunction func(ctx context.Context, app *app.App, params SlackParams, endpointContext *Context) error

// SlackParams maps the main parameters used by a slack bot endpoint
type SlackParams struct {
	TriggerID string
	ChannelID string
	UserID    string
}

// SlackEndpoint abstract some processing for most slack bot endpoints
type SlackEndpoint struct {
	App      *app.App
	Name     string
	Function SlackEndpointFunction
}

// NewSlackEndpoint creates a new endpoint abstraction for slack requests
func NewSlackEndpoint(app *app.App, name string, function SlackEndpointFunction, errorHandler ErrorHandler) *Endpoint {
	slackEndpoint := SlackEndpoint{
		App:      app,
		Name:     name,
		Function: function,
	}

	return &Endpoint{
		App:          app,
		Name:         name,
		ErrorHandler: errorHandler,
		Function:     slackEndpoint.processSlackFunction,
	}
}

func (se SlackEndpoint) processSlackFunction(ctx context.Context, app *app.App, endpointContext *Context) error {
	var slackParams SlackParams

	endpointContext.ReadForm().
		Read(&slackParams.TriggerID, "trigger_id").
		Read(&slackParams.ChannelID, "channel_id").
		Read(&slackParams.UserID, "user_id")

	return se.Function(ctx, app, slackParams, endpointContext)
}

// NewDefaultSlackErrorHandler creates the default error handler for slack endpoints
func NewDefaultSlackErrorHandler() ErrorHandler {
	return defaultSlackErrorHandler
}

func defaultSlackErrorHandler(ctx context.Context, endpoint *Endpoint, endpointContext *Context, err error) {
	if err != nil {
		endpoint.App.Logger.Error(
			ctx,
			log.Trace(),
			log.Reason(fmt.Sprintf("commands.%s", endpoint.Name)),
			log.NewValue("error", err),
		)

		endpointContext.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	endpointContext.Status(http.StatusOK)
}
