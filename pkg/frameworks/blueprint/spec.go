package blueprint

import (
	"strings"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"github.com/blueprint-uservices/blueprint/plugins/memcached"
	"github.com/blueprint-uservices/blueprint/plugins/mongodb"
	"github.com/blueprint-uservices/blueprint/plugins/mysql"
	"github.com/blueprint-uservices/blueprint/plugins/rabbitmq"
	"github.com/blueprint-uservices/blueprint/plugins/redis"
	"github.com/blueprint-uservices/blueprint/plugins/workflow"
	"github.com/blueprint-uservices/blueprint/plugins/workflow/workflowspec"

	"analyzer/pkg/datastores"
	"analyzer/pkg/frameworks"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

func LoadWiring(appName string) ([]*frameworks.ServiceInfo, []datastores.DatabaseInstance, []string) {
	spec := utils.LoadAppSpec(appName)
	servicesSpec, databasesNodes, frontends := BuildAndInspectIR(appName, spec)
	servicesInfo := buildBlueprintServicesInfo(servicesSpec)
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

func buildBlueprintServicesInfo(appSpecs map[*workflowspec.Service][]golang.Service) []*frameworks.ServiceInfo {
	var services []*frameworks.ServiceInfo
	for spec, serviceArgs := range appSpecs {
		serviceInfo := &frameworks.ServiceInfo{
			Name:            spec.Iface.Name,
			Package:         spec.Iface.File.Package.ShortName,
			PackagePath:     spec.Iface.File.Package.Name,
			Filepath:        spec.Iface.File.Name,
			ConstructorName: spec.Constructor.Name,
			ServiceArgs:     []string{"context"}, // args in spec do not count with the context at index 0 so we insert a dummy value now
		}
		for _, method := range spec.Iface.Ast.Methods.List {
			serviceInfo.Methods = append(serviceInfo.Methods, method.Names[0].Name)
		}
		for _, arg := range serviceArgs {
			if workflowClient, ok := arg.(*workflow.WorkflowClient); ok {
				serviceInfo.Edges = append(serviceInfo.Edges, getUniqueName(workflowClient.ServiceType))
			}
			serviceInfo.ServiceArgs = append(serviceInfo.ServiceArgs, getUniqueName(arg.Name()))
		}
		services = append(services, serviceInfo)
	}
	return services
}

func buildDatabasesInstances(databases map[string]ir.IRNode) []datastores.DatabaseInstance {
	var dbInstances []datastores.DatabaseInstance
	for name, node := range databases {
		name = getUniqueName(name)
		switch node.(type) {
		case *redis.RedisGoClient:
			dbInstances = append(dbInstances, &CacheInstance{
				BlueprintDatabaseInstance: BlueprintDatabaseInstance{
					Name: name,
					//FIXME, we can have many replicas
					Datastore: &datastores.Datastore{
						Type:   datastores.Cache,
						Kind:   datastores.Redis,
						Name:   name,
						Schema: &datastores.Schema{},
					},
				},
			})
		case *memcached.MemcachedGoClient:
			dbInstances = append(dbInstances, &CacheInstance{
				BlueprintDatabaseInstance: BlueprintDatabaseInstance{
					Name: name,
					//FIXME, we can have many replicas
					Datastore: &datastores.Datastore{
						Type:   datastores.Cache,
						Kind:   datastores.Memcached,
						Name:   name,
						Schema: &datastores.Schema{},
					},
				},
			})
		case *rabbitmq.RabbitmqGoClient:
			dbInstances = append(dbInstances, &QueueInstance{
				BlueprintDatabaseInstance: BlueprintDatabaseInstance{
					Name: name,
					//FIXME, we can have many replicas
					Datastore: &datastores.Datastore{
						Type:   datastores.Queue,
						Kind:   datastores.RabbitMQ,
						Name:   name,
						Schema: &datastores.Schema{},
					},
				},
			})
		case *mongodb.MongoDBGoClient:
			dbInstances = append(dbInstances, &NoSQLInstance{
				BlueprintDatabaseInstance: BlueprintDatabaseInstance{
					Name: name,
					//FIXME, we can have many replicas
					Datastore: &datastores.Datastore{
						Type:   datastores.NoSQL,
						Kind:   datastores.MongoDB,
						Name:   name,
						Schema: &datastores.Schema{},
					},
				},
			})
		case *mysql.MySQLDBGoClient:
			dbInstances = append(dbInstances, &RelationalDBInstance{
				BlueprintDatabaseInstance: BlueprintDatabaseInstance{
					Name: name,
					//FIXME, we can have many replicas
					Datastore: &datastores.Datastore{
						Type:   datastores.SQL,
						Kind:   datastores.MySQL,
						Name:   name,
						Schema: &datastores.Schema{},
					},
				},
			})
		default:
			logger.Logger.Fatalf("unknown type for database instance: %s (type = %s)", name, utils.GetType(node))
			continue
		}
	}
	return dbInstances
}
