module github.com/blueprint-uservices/blueprint/examples/dsb_hotel2/init

go 1.22.4

require (
	github.com/blueprint-uservices/blueprint/examples/dsb_hotel2/workflow v0.0.0
	github.com/blueprint-uservices/blueprint/runtime v0.0.0-20240405152959-f078915d2306
)

require (
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
)

replace github.com/blueprint-uservices/blueprint/examples/dsb_hotel2/workflow => ../../../../../../blueprint/examples/dsb_hotel2/workflow
