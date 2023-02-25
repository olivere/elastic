module github.com/olivere/elastic/v7/recipes

go 1.17

require (
	github.com/aws/aws-sdk-go v1.43.21
	github.com/google/uuid v1.3.0
	github.com/olivere/elastic/v7 v7.0.0-00010101000000-000000000000
	github.com/olivere/env v1.1.0
	github.com/smartystreets/go-aws-auth v0.0.0-20180515143844-0c1422d1fdb9
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	go.opentelemetry.io/otel v1.5.0
	go.opentelemetry.io/otel/exporters/jaeger v1.5.0
	go.opentelemetry.io/otel/sdk v1.5.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require (
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/smartystreets/gunit v1.4.3 // indirect
	go.opentelemetry.io/otel/trace v1.5.0 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace github.com/olivere/elastic/v7 => ../../elastic
