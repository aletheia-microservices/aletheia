package specs

import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/gotests"
	"github.com/blueprint-uservices/blueprint/plugins/mongodb"
	"github.com/blueprint-uservices/blueprint/plugins/rabbitmq"
	"github.com/blueprint-uservices/blueprint/plugins/workflow"
)

var Docker = cmdbuilder.SpecOption{
	Name:        "docker",
	Description: "Deploys each service in a separate container with thrift, and uses mongodb as NoSQL database backends.",
	Build:       makeDockerSpec,
}

// Create a basic social network wiring spec.
// Returns the names of the nodes to instantiate or an error.
func makeDockerSpec(spec wiring.WiringSpec) ([]string, error) {
	var containers []string
	var allServices []string

	posts_db := mongodb.Container(spec, "posts_db")
	analytics_db := mongodb.Container(spec, "analytics_db")
	notifications_queue := rabbitmq.Container(spec, "notifications_queue", "notifications_queue")
	analytics_queue := rabbitmq.Container(spec, "analytics_queue", "analytics_queue")

	allServices = append(allServices, posts_db)
	allServices = append(allServices, analytics_db)
	allServices = append(allServices, notifications_queue)
	allServices = append(allServices, analytics_queue)

	analytics_service := workflow.Service[postnotification_simple.AnalyticsService](spec, "analytics_service", analytics_db, analytics_queue)
	analytics_service_ctr := applyDockerDefaults(spec, analytics_service, "analytics_service_proc", "analytics_service_container")
	containers = append(containers, analytics_service_ctr)
	allServices = append(allServices, "analytics_service")

	storage_service := workflow.Service[postnotification_simple.StorageService](spec, "storage_service", posts_db, analytics_queue)
	storage_service_ctr := applyDockerDefaults(spec, storage_service, "storage_service_proc", "storage_service_container")
	containers = append(containers, storage_service_ctr)
	allServices = append(allServices, "storage_service")

	notify_service := workflow.Service[postnotification_simple.NotifyService](spec, "notify_service", storage_service, notifications_queue)
	notify_service_ctr := applyDockerQueueHandlerDefaults(spec, notify_service, "notify_service_proc", "notify_service_container")
	containers = append(containers, notify_service_ctr)

	upload_service := workflow.Service[postnotification_simple.UploadService](spec, "upload_service", storage_service, analytics_service, notifications_queue)
	upload_service_ctr := applyHTTPDefaults(spec, upload_service, "upload_service_proc", "upload_service_container")
	containers = append(containers, upload_service_ctr)
	allServices = append(allServices, "upload_service")

	tests := gotests.Test(spec, allServices...)
	containers = append(containers, tests)

	return containers, nil
}
