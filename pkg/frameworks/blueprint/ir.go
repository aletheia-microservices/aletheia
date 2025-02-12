package blueprint

import (
	"reflect"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/blueprint/logging"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/address"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/namespaceutil"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"github.com/blueprint-uservices/blueprint/plugins/goproc"
	"github.com/blueprint-uservices/blueprint/plugins/http"
	"github.com/blueprint-uservices/blueprint/plugins/linuxcontainer"
	"github.com/blueprint-uservices/blueprint/plugins/memcached"
	"github.com/blueprint-uservices/blueprint/plugins/mongodb"
	"github.com/blueprint-uservices/blueprint/plugins/mysql"
	"github.com/blueprint-uservices/blueprint/plugins/rabbitmq"
	"github.com/blueprint-uservices/blueprint/plugins/redis"
	"github.com/blueprint-uservices/blueprint/plugins/workflow"
	"github.com/blueprint-uservices/blueprint/plugins/workflow/workflowspec"

	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

func BuildAndInspectIR(name string, spec cmdbuilder.SpecOption) (map[*workflowspec.Service][]golang.Service, map[string]ir.IRNode, []string) {
	builder := buildIR(name, spec)
	return inspectIR(builder)
}

func buildIR(name string, spec cmdbuilder.SpecOption) *cmdbuilder.CmdBuilder {
	builder := &cmdbuilder.CmdBuilder{
		Name:      name,
		Registry:  map[string]cmdbuilder.SpecOption{},
		Spec:      spec,
		OutputDir: "build",
	}
	logging.DisableCompilerLogging()
	builder.Registry[spec.Name] = spec
	builder.Wiring = wiring.NewWiringSpec(builder.Name)
	if builder.Wiring == nil {
		logger.Logger.Fatal("error creating new wiring spec")
		return nil
	}
	nodesToBuild, err := builder.Spec.Build(builder.Wiring)
	if err != nil {
		logger.Logger.Fatal("error building wiring spec")
		return nil
	}

	builder.IR, err = builder.Wiring.BuildIR(nodesToBuild...)
	if err != nil {
		logger.Logger.Fatal("error building IR")
		return nil
	}
	return builder
}

func inspectIR(builder *cmdbuilder.CmdBuilder) (map[*workflowspec.Service][]golang.Service, map[string]ir.IRNode, []string) {
	services := make(map[*workflowspec.Service][]golang.Service)
	databases := make(map[string]ir.IRNode)
	var frontends []string
	logger.Logger.Infof("[IR] inspecting ir %v", builder.IR)
	logger.Logger.Trace()
	for _, node := range builder.IR.Children {
		if n, ok := node.(*address.Address[*http.GolangHttpServer]); ok {
			if httpService, ok := n.GetDestination().(*http.GolangHttpServer); ok {
				if workflowHandler, ok := httpService.Wrapped.(*workflow.WorkflowHandler); ok {
					frontends = append(frontends, workflowHandler.ServiceInfo.Iface.Name)
				}
			}
		}
		if n, ok := node.(namespaceutil.IRNamespace); ok {
			if nn, ok := n.(ir.IRNode); ok {
				if nnn, ok := nn.(*linuxcontainer.Container); ok {
					logger.Logger.Trace("--------------------------------------------")
					logger.Logger.Trace()
					logger.Logger.Trace(nnn.String())
					logger.Logger.Trace()
					for _, child := range nnn.Nodes {
						if nnnn, ok := child.(*goproc.Process); ok {
							for _, child := range nnnn.Edges {
								t := reflect.TypeOf(child).Elem().Name()
								logger.Logger.Tracef("[IR EDGE] got edge %s with type %s", child.Name(), t)
							}
							for _, child := range nnnn.Nodes {
								if redisClient, ok := child.(*redis.RedisGoClient); ok {
									logger.Logger.Tracef("[IR NODE] [redis.RedisGoClient] got node %s", redisClient.Name())
								} else if rabbitClient, ok := child.(*rabbitmq.RabbitmqGoClient); ok {
									logger.Logger.Tracef("[IR NODE] [rabbitmq.RabbitmqGoClient] got node %s", rabbitClient.Name())
								} else if workflowHandler, ok := child.(*workflow.WorkflowHandler); ok {
									logger.Logger.Tracef("[IR NODE] [workflow.WorkflowHandler] got node %s (service_type = %v)", workflowHandler.Name(), workflowHandler.ServiceType)

									if workflowHandler.ServiceType == "Runnable" {
										logger.Logger.Fatalf("[IR NODE] found Runnable service type for service (%s) -- cannot analyze application", workflowHandler.Name())
									}

									// make sure that services that do not have any other dependencies are also included
									services[workflowHandler.ServiceInfo] = nil
									for _, arg := range workflowHandler.Args {
										if redisClient, ok := arg.(*redis.RedisGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], redisClient)
											databases[redisClient.Name()] = redisClient
											logger.Logger.Tracef("[IR HANDLER ARG] [redis.RedisGoClient] got node %s", redisClient.Name())
										} else if memcachedClient, ok := arg.(*memcached.MemcachedGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], memcachedClient)
											databases[memcachedClient.Name()] = memcachedClient
											logger.Logger.Tracef("[IR HANDLER ARG] [memcached.MemcachedGoClient] got node %s", memcachedClient.Name())
										} else if rabbitClient, ok := arg.(*rabbitmq.RabbitmqGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], rabbitClient)
											databases[rabbitClient.Name()] = rabbitClient
											logger.Logger.Tracef("[IR HANDLER ARG] [rabbitmq.RabbitmqGoClient] got node %s", rabbitClient.Name())
										} else if mongoDbClient, ok := arg.(*mongodb.MongoDBGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], mongoDbClient)
											databases[mongoDbClient.Name()] = mongoDbClient
											logger.Logger.Tracef("[IR HANDLER ARG] [mongodb.MongoDBGoClient] got node %s", mongoDbClient.Name())
										} else if mysqlClient, ok := arg.(*mysql.MySQLDBGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], mysqlClient)
											databases[mysqlClient.Name()] = mysqlClient
											logger.Logger.Tracef("[IR HANDLER ARG] [mongodb.MongoDBGoClient] got node %s", mysqlClient.Name())
										} else if workflowClient, ok := arg.(*workflow.WorkflowClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], workflowClient)
											logger.Logger.Tracef("[IR HANDLER ARG] [workflow.WorkflowClient] got node %s (service_type = %v)", workflowClient.Name(), workflowClient.ServiceType)
										}
									}
								}
							}
						} else {
							logger.Logger.Tracef("unknown node type: %s", utils.GetType(child))
						}
					}
					logger.Logger.Trace()
				}
			}
		} else if redisContainer, ok := node.(*redis.RedisContainer); ok {
			logger.Logger.Tracef("[IR INFO] ignoring redis.RedisContainer for node %s, interface %s", redisContainer.Name(), redisContainer.Iface)
		} else if memachedContainer, ok := node.(*memcached.MemcachedContainer); ok {
			logger.Logger.Tracef("[IR INFO] ignoring memcached.MemcachedContainer for node %s, interface %s", memachedContainer.Name(), memachedContainer.Iface)
		} else if rabbitContainer, ok := node.(*rabbitmq.RabbitmqContainer); ok {
			logger.Logger.Tracef("[IR INFO] ignoring rabbitmq.RabbitmqContainer for node %s, interface %s", rabbitContainer.Name(), rabbitContainer.Iface)
		} else if mongoDbContainer, ok := node.(*mongodb.MongoDBContainer); ok {
			logger.Logger.Tracef("[IR INFO] ignoring mongodb.MongoDBContainer for node %s, interface %s", mongoDbContainer.Name(), mongoDbContainer.Iface)
		} else if mysqlContainer, ok := node.(*mysql.MySQLDBContainer); ok {
			logger.Logger.Tracef("[IR INFO] ignoring mysql.MySQLDBContainer for node %s, interface %s", mysqlContainer.Name(), mysqlContainer.Iface)
		}
	}

	logger.Logger.Trace()
	for key, value := range services {
		logger.Logger.Tracef("[IR SERVICE] inspecting service %s", key.Iface.Name)
		for _, arg := range value {
			if workflowClient, ok := arg.(*workflow.WorkflowClient); ok {
				logger.Logger.Tracef("[IR SERVICE] \t\t[workflow] %s", workflowClient.ServiceType)
			} 
		}
	}
	logger.Logger.Trace()
	logger.Logger.Tracef("[IR DATASTORE] inspecting databases")
	for _, value := range databases {
		if rabbitClient, ok := value.(*rabbitmq.RabbitmqGoClient); ok {
			logger.Logger.Tracef("[IR DATASTORE] \t\t[rabbitmq] %s", rabbitClient.Name())
		} else if redisClient, ok := value.(*redis.RedisGoClient); ok {
			logger.Logger.Tracef("[IR DATASTORE] \t\t[redis] %s", redisClient.Name())
		} else if memcachedClient, ok := value.(*memcached.MemcachedGoClient); ok {
			logger.Logger.Tracef("[IR DATASTORE] \t\t[memcached] %s", memcachedClient.Name())
		} else if mongodbClient, ok := value.(*mongodb.MongoDBGoClient); ok {
			logger.Logger.Tracef("[IR DATASTORE] \t\t[mongodb] %s", mongodbClient.Name())
		} else if mysqlClient, ok := value.(*mysql.MySQLDBGoClient); ok {
			logger.Logger.Tracef("[IR DATASTORE] \t\t[mysqlClient] %s", mysqlClient.Name())
		}
	}
	return services, databases, frontends
}
