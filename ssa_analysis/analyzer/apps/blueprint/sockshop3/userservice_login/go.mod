module github.com/blueprint-uservices/blueprint/examples/sockshop3/init

go 1.22.4

require (
	github.com/blueprint-uservices/blueprint/examples/sockshop3/workflow v0.0.0
	github.com/blueprint-uservices/blueprint/runtime v0.0.0-20250729202253-a8f505263256
)

require (
	github.com/Rhymond/go-money v1.0.15 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.mongodb.org/mongo-driver v1.17.4 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
)

replace github.com/blueprint-uservices/blueprint/examples/sockshop3/workflow => ../../../../../../blueprint/examples/sockshop3/workflow
