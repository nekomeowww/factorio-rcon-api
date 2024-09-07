package http

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TraceparentWrapper struct {
	next  http.Handler
	props propagation.TextMapPropagator
}

func NewTraceparentWrapper(next http.Handler) *TraceparentWrapper {
	return &TraceparentWrapper{
		next:  next,
		props: otel.GetTextMapPropagator(),
	}
}

func (t *TraceparentWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.props.Inject(r.Context(), propagation.HeaderCarrier(w.Header()))
	t.next.ServeHTTP(w, r)
}
