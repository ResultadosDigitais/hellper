package endpoint

import "net/http"

// Context is responsible for providing access to useful functions to your endpoint function. It
// encapsulates most of the work.
type Context struct {
	writer  http.ResponseWriter
	request *http.Request
}

// Status writes the HTTP status to the response
func (ec *Context) Status(httpStatus int) {
	ec.writer.WriteHeader(httpStatus)
}

// Error writes the error message and the HTTP status to the response
func (ec *Context) Error(errorMessage string, httpStatus int) {
	http.Error(ec.writer, errorMessage, httpStatus)
}

// Header adds or replaces a header value in the response
func (ec *Context) Header(name, value string) {
	ec.writer.Header().Set(name, value)
}

// ReadForm enables an endpoint function to read variables from the request
func (ec *Context) ReadForm() *FormReader {
	return &FormReader{request: ec.request}
}

// FormReader is responsible for encapsulating the variable reading for the endpoint function
type FormReader struct {
	request *http.Request
}

// Read a variable from the form and stores it inside the target variable
func (fr *FormReader) Read(target *string, formVariableName string) *FormReader {
	*target = fr.request.FormValue(formVariableName)
	return fr
}
