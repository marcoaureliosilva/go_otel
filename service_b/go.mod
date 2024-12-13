module service_b

go 1.22.0

toolchain go1.22.3

require (
	github.com/openzipkin/zipkin-go v0.4.3
	go.opentelemetry.io/contrib v1.33.0 // Certifique-se de que você tenha essas dependências
	go.opentelemetry.io/otel v1.33.0 // ou a versão mais recente
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/sdk v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)
