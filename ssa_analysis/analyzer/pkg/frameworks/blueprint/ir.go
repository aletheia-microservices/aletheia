package blueprint

import (
	"fmt"
	"log"
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
	logging.EnableCompilerLogging()
	builder.Registry[spec.Name] = spec
	builder.Wiring = wiring.NewWiringSpec(builder.Name)
	if builder.Wiring == nil {
		log.Fatal("error creating new wiring spec")
		return nil
	}
	nodesToBuild, err := builder.Spec.Build(builder.Wiring)
	if err != nil {
		log.Fatal("error building wiring spec")
		return nil
	}

	builder.IR, err = builder.Wiring.BuildIR(nodesToBuild...)
	if err != nil {
		log.Fatal("error building IR")
		return nil
	}
	return builder
}

func inspectIR(builder *cmdbuilder.CmdBuilder) (map[*workflowspec.Service][]golang.Service, map[string]ir.IRNode, []string) {
	services := make(map[*workflowspec.Service][]golang.Service)
	databases := make(map[string]ir.IRNode)
	var frontends []string
	fmt.Printf("[IR] inspecting ir %v\n", builder.IR)
	fmt.Println()
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
					fmt.Println("--------------------------------------------")
					fmt.Println()
					fmt.Println(nnn.String())
					fmt.Println()
					for _, child := range nnn.Nodes {
						if nnnn, ok := child.(*goproc.Process); ok {
							for _, child := range nnnn.Edges {
								t := reflect.TypeOf(child).Elem().Name()
								fmt.Printf("[IR EDGE] got edge %s with type %s\n", child.Name(), t)
							}
							for _, child := range nnnn.Nodes {
								if redisClient, ok := child.(*redis.RedisGoClient); ok {
									fmt.Printf("[IR NODE] [redis.RedisGoClient] got node %s\n", redisClient.Name())
								} else if rabbitClient, ok := child.(*rabbitmq.RabbitmqGoClient); ok {
									fmt.Printf("[IR NODE] [rabbitmq.RabbitmqGoClient] got node %s\n", rabbitClient.Name())
								} else if workflowHandler, ok := child.(*workflow.WorkflowHandler); ok {
									fmt.Printf("[IR NODE] [workflow.WorkflowHandler] got node %s (service_type = %v)\n", workflowHandler.Name(), workflowHandler.ServiceType)

									if workflowHandler.ServiceType == "Runnable" {
										log.Fatalf("[IR NODE] found Runnable service type for service (%s) -- cannot analyze application", workflowHandler.Name())
									}

									// make sure that services that do not have any other dependencies are also included
									services[workflowHandler.ServiceInfo] = nil
									for _, arg := range workflowHandler.Args {
										if redisClient, ok := arg.(*redis.RedisGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], redisClient)
											databases[redisClient.Name()] = redisClient
											fmt.Printf("[IR HANDLER ARG] [redis.RedisGoClient] got node %s\n", redisClient.Name())
										} else if memcachedClient, ok := arg.(*memcached.MemcachedGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], memcachedClient)
											databases[memcachedClient.Name()] = memcachedClient
											fmt.Printf("[IR HANDLER ARG] [memcached.MemcachedGoClient] got node %s\n", memcachedClient.Name())
										} else if rabbitClient, ok := arg.(*rabbitmq.RabbitmqGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], rabbitClient)
											databases[rabbitClient.Name()] = rabbitClient
											fmt.Printf("[IR HANDLER ARG] [rabbitmq.RabbitmqGoClient] got node %s\n", rabbitClient.Name())
										} else if mongoDbClient, ok := arg.(*mongodb.MongoDBGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], mongoDbClient)
											databases[mongoDbClient.Name()] = mongoDbClient
											fmt.Printf("[IR HANDLER ARG] [mongodb.MongoDBGoClient] got node %s\n", mongoDbClient.Name())
										} else if mysqlClient, ok := arg.(*mysql.MySQLDBGoClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], mysqlClient)
											databases[mysqlClient.Name()] = mysqlClient
											fmt.Printf("[IR HANDLER ARG] [mongodb.MongoDBGoClient] got node %s\n", mysqlClient.Name())
										} else if workflowClient, ok := arg.(*workflow.WorkflowClient); ok {
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], workflowClient)
											fmt.Printf("[IR HANDLER ARG] [workflow.WorkflowClient] got node %s (service_type = %v)\n", workflowClient.Name(), workflowClient.ServiceType)
										}
									}
								}
							}
						} else {
							fmt.Printf("unknown node type: [%T] %v\n", child, child)
						}
					}
					fmt.Println()
				}
			}
		} else if redisContainer, ok := node.(*redis.RedisContainer); ok {
			fmt.Printf("[IR INFO] ignoring redis.RedisContainer for node %s, interface %s\n", redisContainer.Name(), redisContainer.Iface)
		} else if memachedContainer, ok := node.(*memcached.MemcachedContainer); ok {
			fmt.Printf("[IR INFO] ignoring memcached.MemcachedContainer for node %s, interface %s\n", memachedContainer.Name(), memachedContainer.Iface)
		} else if rabbitContainer, ok := node.(*rabbitmq.RabbitmqContainer); ok {
			fmt.Printf("[IR INFO] ignoring rabbitmq.RabbitmqContainer for node %s, interface %s\n", rabbitContainer.Name(), rabbitContainer.Iface)
		} else if mongoDbContainer, ok := node.(*mongodb.MongoDBContainer); ok {
			fmt.Printf("[IR INFO] ignoring mongodb.MongoDBContainer for node %s, interface %s\n", mongoDbContainer.Name(), mongoDbContainer.Iface)
		} else if mysqlContainer, ok := node.(*mysql.MySQLDBContainer); ok {
			fmt.Printf("[IR INFO] ignoring mysql.MySQLDBContainer for node %s, interface %s\n", mysqlContainer.Name(), mysqlContainer.Iface)
		}
	}

	fmt.Println()
	for key, value := range services {
		fmt.Printf("[IR SERVICE] inspecting service %s\n", key.Iface.Name)
		for _, arg := range value {
			if workflowClient, ok := arg.(*workflow.WorkflowClient); ok {
				fmt.Printf("[IR SERVICE] \t\t[workflow] %s\n", workflowClient.ServiceType)
			}
		}
	}
	fmt.Println()
	fmt.Printf("[IR DATASTORE] inspecting databases\n")
	for _, value := range databases {
		if rabbitClient, ok := value.(*rabbitmq.RabbitmqGoClient); ok {
			fmt.Printf("[IR DATASTORE] \t\t[rabbitmq] %s\n", rabbitClient.Name())
		} else if redisClient, ok := value.(*redis.RedisGoClient); ok {
			fmt.Printf("[IR DATASTORE] \t\t[redis] %s\n", redisClient.Name())
		} else if memcachedClient, ok := value.(*memcached.MemcachedGoClient); ok {
			fmt.Printf("[IR DATASTORE] \t\t[memcached] %s\n", memcachedClient.Name())
		} else if mongodbClient, ok := value.(*mongodb.MongoDBGoClient); ok {
			fmt.Printf("[IR DATASTORE] \t\t[mongodb] %s\n", mongodbClient.Name())
		} else if mysqlClient, ok := value.(*mysql.MySQLDBGoClient); ok {
			fmt.Printf("[IR DATASTORE] \t\t[mysqlClient] %s\n", mysqlClient.Name())
		}
	}
	return services, databases, frontends
}
