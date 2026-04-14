package blueprint

import (
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
	"github.com/sirupsen/logrus"
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
	builder.Registry[spec.Name] = spec
	builder.Wiring = wiring.NewWiringSpec(builder.Name)
	if builder.Wiring == nil {
		logrus.Fatalf("error creating new wiring spec")
		return nil
	}
	nodesToBuild, err := builder.Spec.Build(builder.Wiring)
	if err != nil {
		logrus.Fatalf("error building wiring spec")
		return nil
	}

	builder.IR, err = builder.Wiring.BuildIR(nodesToBuild...)
	if err != nil {
		logrus.Fatalf("error building IR")
		return nil
	}
	return builder
}

func inspectIR(builder *cmdbuilder.CmdBuilder) (map[*workflowspec.Service][]golang.Service, map[string]ir.IRNode, map[*workflowspec.Service][]ir.IRNode, []string) {
	services := make(map[*workflowspec.Service][]golang.Service)
	databases := make(map[string]ir.IRNode)
	args := make(map[*workflowspec.Service][]ir.IRNode)
	var frontends []string
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
					for _, child := range nnn.Nodes {
						if nnnn, ok := child.(*goproc.Process); ok {
							for _, child := range nnnn.Nodes {
								if workflowHandler, ok := child.(*workflow.WorkflowHandler); ok {
									if workflowHandler.ServiceType == "Runnable" {
										logrus.Fatalf("[IR NODE] found Runnable service type for service (%s) -- cannot analyze application", workflowHandler.Name())
									}

									// make sure that services that do not have any other dependencies are also included
									services[workflowHandler.ServiceInfo] = nil

									for _, arg := range workflowHandler.Args {
										switch t := arg.(type) {
										case *redis.RedisGoClient, *memcached.MemcachedGoClient, *rabbitmq.RabbitmqGoClient, *mongodb.MongoDBGoClient, *mysql.MySQLDBGoClient:
											databases[arg.Name()] = arg
										case *workflow.WorkflowClient:
											services[workflowHandler.ServiceInfo] = append(services[workflowHandler.ServiceInfo], t)
										case *ir.IRValue:
											// nothing to do
										default:
											logrus.Fatalf("[IR HANDLER ARG] unknown arg type [%T]: %v\n", arg, arg)
										}
										args[workflowHandler.ServiceInfo] = append(args[workflowHandler.ServiceInfo], arg)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return services, databases, args, frontends
}
