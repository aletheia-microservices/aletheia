package app

import (
	"encoding/json"
	"slices"

	"analyzer/pkg/datastores"
	"analyzer/pkg/frameworks"
	"analyzer/pkg/logger"
	"analyzer/pkg/service"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

type App struct {
	Name              string
	PackagePath       string
	Services          map[string]*service.Service
	Databases         map[string]datastores.DatabaseInstance
	AppPackages       map[string]*types.Package // key is package name (FIXME: should be path actually)
	BlueprintPackages map[string]*types.Package // key is package name (FIXME: should be path actually)
	ExternalPackages  map[string]*types.Package // key is package name (FIXME: should be path actually)
	TaintedVariables  map[string][]objects.Object
	ServiceTypes      map[string]*gotypes.ServiceType
	Dataflows         []*types.Dataflow
}

func (app *App) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name      string                                 `json:"name"`
		Services  map[string]*service.Service            `json:"services"`
		Databases map[string]datastores.DatabaseInstance `json:"datastores"`
	}{
		Name:      app.Name,
		Services:  app.Services,
		Databases: app.Databases,
	})
}

func (app *App) ResetAllDataflows() {
	var dataflowsToKeep []*types.Dataflow
	for _, df := range app.Dataflows {
		if !df.GetObjectDataflow().IsPermanent() {
			objInfo := df.GetObjectDataflow().Variable.GetVariableInfo()
			objInfo.ResetAllDataflows()
		} else {
			dataflowsToKeep = append(dataflowsToKeep, df)
		}
	}
	app.Dataflows = dataflowsToKeep
}

func (app *App) AddDataflow(objDataflow *objects.ObjectDataflow, dbCall *types.ParsedDatabaseCall) {
	app.Dataflows = append(app.Dataflows, types.NewDataflow(objDataflow, dbCall))
}

func (app *App) GetDataflowForObjectDataflow(objDataflow *objects.ObjectDataflow) *types.Dataflow {
	for _, df := range app.Dataflows {
		if df.HasObjectDataflow(objDataflow) {
			return df
		}
	}
	logger.Logger.Fatalf("[APP] could not find dataflow for object dataflow: %s", objDataflow.String())
	return nil
}

func (app *App) GetServices() map[string]*service.Service {
	return app.Services
}

func (app *App) GetDbInstances() map[string]datastores.DatabaseInstance {
	return app.Databases
}

// GetDatabasesReferencingCurrentExceptQueues discards any queue or cache
func (app *App) GetDatabasesReferencingCurrentExceptQueues(curr *datastores.Datastore) []*datastores.Datastore {
	var refs []*datastores.Datastore
	for _, dbInstance := range app.GetDbInstances() {
		ds := dbInstance.GetDatastore()

		// discard any
		if ds.IsQueue() {
			continue
		}
		if ds != curr && ds.IsReferencingDatastore(curr) {
			refs = append(refs, ds)
		}
	}

	return refs
}

func (app *App) GetServicesUsingDatastore(ds *datastores.Datastore) []*service.Service {
	var svcs []*service.Service
	for _, service := range app.GetServices() {
		if service.HasDatastore(ds) {
			svcs = append(svcs, service)
		}
	}
	return svcs
}

func (app *App) String() string {
	str, _ := app.MarshalJSON()
	return string(str)
}

func (app *App) GetDatastoreInstance(name string) datastores.DatabaseInstance {
	if instance, ok := app.Databases[name]; ok {
		return instance
	}
	logger.Logger.Fatalf("[APP] could not find datastore instance for name (%s)", name)
	return nil
}

func (app *App) GetAllAppPackages() []*types.Package {
	var allPkgs []*types.Package
	for _, pkg := range app.AppPackages {
		allPkgs = append(allPkgs, pkg)
	}
	return allPkgs
}

// this is just an helper to later print tainted dataflow in yaml files
func (app *App) AddTaintedVariableIfNotExists(fieldName string, variable objects.Object) {
	if !slices.Contains(app.TaintedVariables[fieldName], variable) {
		app.TaintedVariables[fieldName] = append(app.TaintedVariables[fieldName], variable)
	}
}

func (app *App) GetAppPackages() map[string]*types.Package {
	return app.AppPackages
}

func (app *App) GetBlueprintPackages() map[string]*types.Package {
	return app.BlueprintPackages
}

func (app *App) AddAppPackage(name string, p *types.Package) {
	app.AppPackages[name] = p
}

func (app *App) AddBlueprintPackage(name string, p *types.Package) {
	app.BlueprintPackages[name] = p
}

func (app *App) AddExternalPackage(name string, p *types.Package) {
	app.ExternalPackages[name] = p
}

func (app *App) GetPackage(pkgName string) *types.Package {
	pkg, ok := app.AppPackages[pkgName]
	if !ok {
		packagesStr := ""
		for k, t := range app.AppPackages {
			packagesStr += "- " + k + ": " + t.String() + "\n"
		}
		logger.Logger.Fatalf("could not find package (%s) in app packages list\n%s", pkgName, packagesStr)
	}
	return pkg
}

func (app *App) GetServiceIfExists(name string) *service.Service {
	if service, ok := app.Services[name]; ok {
		return service
	}
	return nil
}

func (app *App) GetService(name string) *service.Service {
	if service, ok := app.Services[name]; ok {
		return service
	}
	logger.Logger.Fatalf("[APP] could not find service (%s) for app: %s", name, app.String())
	return nil
}

func InitApp(name string, servicesInfo []*frameworks.ServiceInfo) (*App, error) {
	app := &App{
		Name:              name,
		PackagePath:       utils.LoadAppPath(name),
		Services:          make(map[string]*service.Service),
		Databases:         make(map[string]datastores.DatabaseInstance),
		AppPackages:       make(map[string]*types.Package),
		BlueprintPackages: make(map[string]*types.Package),
		ExternalPackages:  make(map[string]*types.Package),
		TaintedVariables:  make(map[string][]objects.Object),
		ServiceTypes:      make(map[string]*gotypes.ServiceType),
	}

	for _, info := range servicesInfo {
		app.ServiceTypes[info.Name] = &gotypes.ServiceType{
			Name:        info.Name,
			PackagePath: info.PackagePath,
		}
	}
	logger.Logger.Infof("[APP INIT] loading app %s", app.PackagePath)
	return app, nil
}
