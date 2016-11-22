package ghost

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// A ghost router is actually backed by a gorilla mux, but we add some helpers
// to make building the API a little easier.
type Router struct {
	Mux *mux.Router
}

func NewRouter() *Router {
	return &Router{mux.NewRouter()}
}

// Add the given path to the router, and return a route builder to begin
// constructing the additional properties of the router.
func (r *Router) AddRoute(path string) *RouteBuilder {
	builder := &RouteBuilder{
		path:      path,
		extenders: []Extender{},
		model:     &NullModel{},
		validator: &NullValidator{},
		processor: &NullProcessor{},
		writer:    &NullWriter{},
	}
	route := r.Mux.Handle(path, builder)
	builder.route = route

	return builder
}

type RouteBuilder struct {
	path      string
	route     *mux.Route
	extenders []Extender
	model     RequestModel
	validator Validator
	processor Processor
	writer    Writer
}

// Define the HTTP methods (verbs) that can be used on this route.
func (b *RouteBuilder) Methods(methods ...string) *RouteBuilder {
	b.route.Methods(methods...)
	return b
}

// Add an Extender that extends the request context coming through to the
// HTTP handler.
func (b *RouteBuilder) Extender(e Extender) *RouteBuilder {
	b.extenders = append(b.extenders, e)
	return b
}

// Add a model which is simply a representation of the client input in a Go
// data structure.
func (b *RouteBuilder) Model(m RequestModel) *RouteBuilder {
	b.model = m
	return b
}

// Add a validator, which validates input coming from the HTTP client before
// it gets operated on by a processor or writer.
func (b *RouteBuilder) Validator(v Validator) *RouteBuilder {
	b.validator = v
	return b
}

// A processor operates on the input sent from the client, manipulates, extends
// it, whatever. Ultimately it returns something new.
func (b *RouteBuilder) Processor(p Processor) *RouteBuilder {
	b.processor = p
	return b
}

// Add a writer which serializes the output from a processor.
func (b *RouteBuilder) Writer(w Writer) *RouteBuilder {
	b.writer = w
	return b
}

// Satisfies the http.Handler interface. The route builder will use the
// provided extensions to create a handler and serve HTTP through the given
// extensions.
func (b *RouteBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade context
	for _, e := range b.extenders {
		r = r.WithContext(e.Extend(r.Context()))
	}

	// load model
	if err := b.model.FromRequest(r); err != nil {
		b.serveError(w, err)
		return
	}

	// run validator
	if err := b.validator.Validate(b.model); err != nil {
		b.serveError(w, err)
		return
	}

	// run processor
	if output, err := b.processor.Process(b.model); err != nil {
		b.serveError(w, err)
		return
	} else if serial, err := b.writer.Serialize(output); err != nil {
		b.serveError(w, err)
		return
	} else {
		fmt.Fprint(w, string(serial))
		return
	}
}

func (b *RouteBuilder) serveError(w http.ResponseWriter, err HttpError) {
	http.Error(w, err.Error(), err.StatusCode())
}

// An implementation of HttpError for simple things.
type GhostError struct {
	err  error
	code int
}

func (e *GhostError) Error() string {
	return e.err.Error()
}

func (e *GhostError) StatusCode() int {
	return e.code
}

func NewHttpError(err error, code int) HttpError {
	return &GhostError{err, code}
}
