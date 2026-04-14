package blueprint

import (
	"strings"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"github.com/blueprint-uservices/blueprint/plugins/memcached"
	"github.com/blueprint-uservices/blueprint/plugins/mongodb"
	"github.com/blueprint-uservices/blueprint/plugins/mysql"
	"github.com/blueprint-uservices/blueprint/plugins/rabbitmq"
	"github.com/blueprint-uservices/blueprint/plugins/redis"
	"github.com/blueprint-uservices/blueprint/plugins/workflow"
	"github.com/blueprint-uservices/blueprint/plugins/workflow/workflowspec"
	"github.com/sirupsen/logrus"

	blueprint_apps "analyzer/pkg/frameworks/blueprint/apps"
	"analyzer/pkg/frameworks/components"
)

func loadAppSpec(app string) cmdbuilder.SpecOption {
	logrus.WithField("app", app).Infof("loading app spec")
	if info, ok := blueprint_apps.APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	logrus.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}

func LoadWiring(appName string, synthetic bool) ([]*components.ServiceInfo, []*components.DatastoreInfo, []string) {
	var spec cmdbuilder.SpecOption
	spec = loadAppSpec(appName)
	servicesSpec, databasesNodes, servicesArgs, frontends := BuildAndInspectIR(appName, spec)
	servicesInfo := buildBlueprintServicesInfo(servicesSpec, servicesArgs)
	databasesInfo := buildDatabasesInstances(databasesNodes)
	return servicesInfo, databasesInfo, frontends
}

func getUniqueName(name string) string {
	// remove .client suffix (e.g. notification_queue.client)
	splits := strings.Split(name, ".")[0]
	if len(splits) > 0 {
		return strings.Split(name, ".")[0]
	}
	return ""
}

func buildBlueprintServicesInfo(appSpecs map[*workflowspec.Service][]golang.Service, servicesArgs map[*workflowspec.Service][]ir.IRNode) []*components.ServiceInfo {
	var services []*components.ServiceInfo
	for serviceSpec, otherServicesLst := range appSpecs {
		serviceInfo := &components.ServiceInfo{
			Name:            serviceSpec.Iface.Name,
			Package:         serviceSpec.Iface.File.Package.ShortName,
			PackagePath:     serviceSpec.Iface.File.Package.Name,
			Filepath:        serviceSpec.Iface.File.Name,
			ConstructorName: serviceSpec.Constructor.Name,
			ServiceArgs:     []string{"context"}, // args in spec do not count with the context at index 0 so we insert a dummy value now
		}
		for _, method := range serviceSpec.Iface.Ast.Methods.List {
			serviceInfo.Methods = append(serviceInfo.Methods, method.Names[0].Name)
		}
		for _, otherService := range otherServicesLst {
			if workflowClient, ok := otherService.(*workflow.WorkflowClient); ok {
				serviceInfo.Edges = append(serviceInfo.Edges, getUniqueName(workflowClient.ServiceType))
			}
		}

		for _, arg := range servicesArgs[serviceSpec] {
			serviceInfo.ServiceArgs = append(serviceInfo.ServiceArgs, getUniqueName(arg.Name()))
		}
		services = append(services, serviceInfo)
	}
	return services
}

func buildDatabasesInstances(databases map[string]ir.IRNode) []*components.DatastoreInfo {
	var dbs []*components.DatastoreInfo
	for name, node := range databases {
		name = getUniqueName(name)
		switch node.(type) {
		case *redis.RedisGoClient:
			dbs = append(dbs, &components.DatastoreInfo{
				Name: name,
				Type: components.DATASTORE_TYPE_CACHE,
				Kind: components.DATASTORE_KIND_REDIS,
			})
		case *memcached.MemcachedGoClient:
			dbs = append(dbs, &components.DatastoreInfo{
				Name: name,
				Type: components.DATASTORE_TYPE_CACHE,
				Kind: components.DATASTORE_KIND_MEMCACHED,
			})
		case *rabbitmq.RabbitmqGoClient:
			dbs = append(dbs, &components.DatastoreInfo{
				Name: name,
				Type: components.DATASTORE_TYPE_QUEUE,
				Kind: components.DATASTORE_KIND_RABBITMQ,
			})
		case *mongodb.MongoDBGoClient:
			dbs = append(dbs, &components.DatastoreInfo{
				Name: name,
				Type: components.DATASTORE_TYPE_NOSQL,
				Kind: components.DATASTORE_KIND_MONGODB,
			})
		case *mysql.MySQLDBGoClient:
			dbs = append(dbs, &components.DatastoreInfo{
				Name: name,
				Type: components.DATASTORE_TYPE_RELATIONALDB,
				Kind: components.DATASTORE_KIND_MYSQL,
			})
		default:
			logrus.Fatalf("unknown type for database instance: %s // NODE = [%T] %v", name, node, node)
			continue
		}
	}
	return dbs
}
