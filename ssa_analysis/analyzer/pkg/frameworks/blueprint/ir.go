package blueprint

import (
	"log"
	/* "reflect" */

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

func BuildAndInspectIR(name string, spec cmdbuilder.SpecOption) (map[*workflowspec.Service][]golang.Service, map[string]ir.IRNode, map[*workflowspec.Service][]ir.IRNode, []string) {
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

func inspectIR(builder *cmdbuilder.CmdBuilder) (map[*workflowspec.Service][]golang.Service, map[string]ir.IRNode, map[*workflowspec.Service][]ir.IRNode, []string) {
	services := make(map[*workflowspec.Service][]golang.Service)
	databases := make(map[string]ir.IRNode)
	args := make(map[*workflowspec.Service][]ir.IRNode)
	var frontends []string
	//EVAL - fmt.Printf("[IR] inspecting ir %v\n", builder.IR)
	//EVAL - fmt.Println()
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
					//EVAL - fmt.Println("--------------------------------------------")
					//EVAL - fmt.Println()
					//EVAL - fmt.Println(nnn.String())
					//EVAL - fmt.Println()
					for _, child := range nnn.Nodes {
						if nnnn, ok := child.(*goproc.Process); ok {
							/* for _, child := range nnnn.Edges { */
								/* t := reflect.TypeOf(child).Elem().Name() */
								//EVAL - fmt.Printf("[IR EDGE] got edge %s with type %s\n", child.Name(), t)
							/* } */
							for _, child := range nnnn.Nodes {
								if workflowHandler, ok := child.(*workflow.WorkflowHandler); ok {
									//EVAL - fmt.Printf("[IR NODE] [workflow.WorkflowHandler] got node %s (service_type = %v)\n", workflowHandler.Name(), workflowHandler.ServiceType)

									if workflowHandler.ServiceType == "Runnable" {
										log.Fatalf("[IR NODE] found Runnable service type for service (%s) -- cannot analyze application", workflowHandler.Name())
									}

									// make sure that services that do not have any other dependencies are also included
									services[workflowHandler.ServiceInfo] = nil

									for _, arg := range workflowHandler.Args {
										//EVAL - fmt.Printf("[IR HANDLER ARG] [%T] got node: %s\n", arg, arg)
										switch t := arg.(type) {
										case *redis.RedisGoClient, *memcached.MemcachedGoClient, *rabbitmq.RabbitmqGoClient, *mongodb.MongoDBGoClient, *mysql.MySQLDBGoClient:
											databases[arg.Name()] = arg
										case *workflow.WorkflowClient:
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], t)
										case *ir.IRValue:
											// nothing to do
										default:
											log.Fatalf("[IR HANDLER ARG] unknown arg type [%T]: %v\n", arg, arg)
										}
										args[workflowHandler.ServiceInfo] = append(args[workflowHandler.ServiceInfo], arg)
									}
								}
							}
						} else {
							//EVAL - fmt.Printf("unknown node type: [%T] %v\n", child, child)
						}
					}
					//EVAL - fmt.Println()
				}
			}
		} /* else if redisContainer, ok := node.(*redis.RedisContainer); ok {
			//EVAL - fmt.Printf("[IR INFO] ignoring redis.RedisContainer for node %s, interface %s\n", redisContainer.Name(), redisContainer.Iface)
		} else if memachedContainer, ok := node.(*memcached.MemcachedContainer); ok {
			//EVAL - fmt.Printf("[IR INFO] ignoring memcached.MemcachedContainer for node %s, interface %s\n", memachedContainer.Name(), memachedContainer.Iface)
		} else if rabbitContainer, ok := node.(*rabbitmq.RabbitmqContainer); ok {
			//EVAL - fmt.Printf("[IR INFO] ignoring rabbitmq.RabbitmqContainer for node %s, interface %s\n", rabbitContainer.Name(), rabbitContainer.Iface)
		} else if mongoDbContainer, ok := node.(*mongodb.MongoDBContainer); ok {
			//EVAL - fmt.Printf("[IR INFO] ignoring mongodb.MongoDBContainer for node %s, interface %s\n", mongoDbContainer.Name(), mongoDbContainer.Iface)
		} else if mysqlContainer, ok := node.(*mysql.MySQLDBContainer); ok {
			//EVAL - fmt.Printf("[IR INFO] ignoring mysql.MySQLDBContainer for node %s, interface %s\n", mysqlContainer.Name(), mysqlContainer.Iface)
		} */
	}
	return services, databases, args, frontends
}
