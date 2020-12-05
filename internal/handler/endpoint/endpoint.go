// Provides an abstraction layer for creating slack bot endpoints
package endpoint

import (
	"context"
	"fmt"
	"hellper/internal/app"
	"hellper/internal/log"
	"net/http"
)

// Endpoint is the configuration about the endpoint to be executed and serves as an HTTP middleware
// for your endpoint.
type Endpoint struct {
	App          *app.App
	Name         string
	Function     func(context.Context, *app.App, *Context) error
	ErrorHandler ErrorHandler
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	r.ParseForm()

	var formValues []log.Value
	for key, value := range r.Form {
		formValues = append(formValues, log.NewValue(key, value))
	}

	e.App.Logger.Debug(
		ctx,
		fmt.Sprintf("endpoint/%s.ServeHTTP Form", e.Name),
		formValues...,
	)

	endpointContext := Context{
		writer:  w,
		request: r,
	}

	err := e.Function(ctx, e.App, &endpointContext)

	if err != nil {
		e.ErrorHandler(ctx, e, &endpointContext, err)
	}
}
