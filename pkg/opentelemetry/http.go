package opentelemetry

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

// HTTPClientTransporter is a convenience function which helps to attach tracing
// functionality to conventional HTTP clients.
func HTTPClientTransporter(rt http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(rt)
}
