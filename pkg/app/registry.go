package app

import (
	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/controlflow"
	"analyzer/pkg/datastores"
	"analyzer/pkg/frameworks"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/lookup"
	"analyzer/pkg/service"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
)

func (app *App) RegisterDatastores(instances []datastores.DatabaseInstance) {
	for _, instance := range instances {
		app.Databases[instance.GetName()] = instance
		logger.Logger.Infof("[APP] registered database instance %s", instance.String())
	}
}

func (app *App) RegisterServices(servicesInfo []*frameworks.ServiceInfo) {
	app.createServiceNodes(servicesInfo)
	app.matchServiceEdges()
	app.buildServiceInfo()
	app.loadFields(servicesInfo)
}

func (app *App) createServiceNodes(servicesInfo []*frameworks.ServiceInfo) {
	logger.Logger.Tracef("[APP] loading #%d specs", len(servicesInfo))
	// services also include blueprint backends
	for _, info := range servicesInfo {
		pkg := app.GetPackage(info.Package)
		file := pkg.GetFile(info.Filepath)

		node := &service.Service{
			Name:                info.Name,
			File:                file,
			Fields:              make(map[string]types.Field),
			Services:            make(map[string]*service.Service),
			Datastores:          make(map[string]datastores.DatabaseInstance),
			ExposedMethods:      make(map[string]*types.ParsedMethod),
			QueueHandlerMethods: make(map[string]*types.ParsedMethod),
			InternalMethods:     make(map[string]*types.ParsedMethod),
			PackageMethods:      make(map[string]*types.ParsedMethod),
			ConstructorName:     info.ConstructorName,
			ImplType: &gotypes.ServiceType{
				PackagePath: pkg.GetName(),
				Name:        info.Name,
			},
			ContextStack: stack.New(),
		}

		// add entries to be later parsed
		for _, name := range info.Methods {
			node.ExposedMethods[name] = nil
		}
		for _, name := range info.Edges {
			node.Services[name] = nil
		}

		app.Services[node.Name] = node
		logger.Logger.Infof("[APP] created service node %s", node.Name)
	}
}

func (app *App) matchServiceEdges() {
	for _, node := range app.Services {
		for edgeName := range node.Services {
			node.Services[edgeName] = app.Services[edgeName]
		}
	}
}

func (app *App) buildServiceInfo() {
	for _, node := range app.Services {
		lookup.ParseFileImports(node.File)
		node.RegisterImplName()
		node.RegisterConstructor()
	}
}

func (app *App) BuildServices() {
	// 1. attach all methods to corresponding services
	for _, node := range app.Services {
		serviceImplementedMethods := node.AttachParsedMethods()
		for _, method := range serviceImplementedMethods {
			controlflow.InitServiceReceiverFieldsForParsedCFG(node, method)
		}
	}
	// 2. attach all package methods that have no service assigned
	for _, node := range app.Services {
		node.AttachAllPackageMethods()
	}
}

func (app *App) loadFields(servicesInfo []*frameworks.ServiceInfo) {
	logger.Logger.Infof("[APP] loading fields for services constructors")
	for _, info := range servicesInfo {
		service := app.Services[info.Name]

		entryBlock := service.GetConstructor().GetParsedCfg().GetEntryParsedBlock()

		for idx, v := range entryBlock.GetVariables() {
			if backendType, ok := v.GetType().(*blueprint.BlueprintBackendType); ok {
				serviceArg := info.GetArgumentAt(idx)
				instance := app.GetDatastoreInstance(serviceArg)
				backendType.SetInstance(instance)
				service.AddDatastore(instance)
			}
		}

		existBlock := service.GetConstructor().GetParsedCfg().GetLastLiveBlock()

		controlflow.ParseMethodCFG(service.GetPackage(), service, service.GetConstructor())
		logger.Logger.Infof("[APP] loading fields for entry block: %s", existBlock.String())
		serviceImplVar := existBlock.GetFirstResult()
		service.SetImplVariableWithType(serviceImplVar)
		service.RegisterFields()

		logger.Logger.Tracef("[APP] registered service node %s with %d service(s) and %d database(s)", service.Name, len(service.Services), len(service.Datastores))
	}
}
