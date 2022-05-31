module github.com/robertscherbarth/opentelemetry-go-example

go 1.18

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.3.0
	github.com/riandyrn/otelchi v0.4.0
	github.com/uptrace/opentelemetry-go-extra/otelzap v0.1.13
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.7.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelutil v0.1.13 // indirect
	go.opentelemetry.io/contrib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)
