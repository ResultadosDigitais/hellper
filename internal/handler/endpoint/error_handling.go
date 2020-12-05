package endpoint

import "context"

// ErrorHandler is the function that will be executed in case of the endpoint executing results in an error
type ErrorHandler func(ctx context.Context, endpoint *Endpoint, endpointContext *Context, err error)
