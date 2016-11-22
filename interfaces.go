package ghost

import (
	"context"
	"net/http"
)

// A type for HTTP errors that extend the base error type. Provides a method
// for status codes to return to the HTTP client.
type HttpError interface {
	error
	StatusCode() int
}

// A context extender is an interface that adds values to a context that
// passes through it. For example, if you want database access inside a context
// you might build a context extender to pass along database connections.
type Extender interface {
	// Extend the context by calling ctx.WithValue()
	Extend(ctx context.Context) context.Context
}

// A request model is a type that wraps user input from a request body. Typically
// this will be a struct representation of some JSON object.
type RequestModel interface {
	// Fill the model with the information from the request. One way to do this
	// is to just use json.Unmarshal(), but you might also want some other
	// information extracted.
	FromRequest(r *http.Request) HttpError
}

// This is a type that validates an InputModel after the model has been parsed
// from the request body.
type Validator interface {
	// Validate the input and return an HTTP error if something is wrong, or
	// nil otherwise.
	Validate(interface{}) HttpError
}

type Processor interface {
	// Process the information or return an HTTP error if something is wrong,
	// or nil otherwise.
	Process(source interface{}) (interface{}, HttpError)
}

// Write the information to the response writer to return the result to the
// HTTP client. Probably this will involve a call to json.Marshal(), unless
// you are using a different serialization method.
type Writer interface {
	Serialize(interface{}) ([]byte, HttpError)
}
