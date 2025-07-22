module github.com/blueprint-uservices/blueprint/examples/postnotification_simple/wiring

go 1.22.4

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/rabbitmq/amqp091-go v1.8.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

require (
	github.com/blueprint-uservices/blueprint/blueprint v0.0.0-20240405152959-f078915d2306
	github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow v0.0.0
)

require (
	github.com/blueprint-uservices/blueprint/plugins v0.0.0-20240405152959-f078915d2306
	github.com/blueprint-uservices/blueprint/runtime v0.0.0-20240405152959-f078915d2306 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/otiai10/copy v1.14.0 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	go.opentelemetry.io/otel v1.25.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.25.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.25.0 // indirect
	go.opentelemetry.io/otel/metric v1.25.0 // indirect
	go.opentelemetry.io/otel/sdk v1.25.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.25.0 // indirect
	go.opentelemetry.io/otel/trace v1.25.0 // indirect
	golang.org/x/exp v0.0.0-20240409090435-93d18d7e34b8 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
)

replace github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow => ../workflow
