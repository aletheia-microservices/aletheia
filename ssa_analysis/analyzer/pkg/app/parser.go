package app

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/app/services"
	"analyzer/pkg/frameworks/blueprint"
)

func (app *App) Init() {
	servicesInfo, datastoresInfo, frontends := blueprint.LoadWiring(app.GetName())

	// parse services
	for _, svcInfo := range servicesInfo {
		name := svcInfo.Name
		constructor := svcInfo.ConstructorName
		pkg := svcInfo.Package
		pkgpath := svcInfo.PackagePath
		path := svcInfo.PackagePath + "." + svcInfo.Name
		impl := svcInfo.Name + "Impl"
		methods := svcInfo.Methods
		args := svcInfo.ServiceArgs

		service := services.NewService(name, impl, pkg, pkgpath, path, constructor, args)
		service.SetMethods(methods...)
		app.AddService(service)
	}
	for _, svcInfo := range servicesInfo {
		service := app.GetServiceByName(svcInfo.Name)
		for _, dep := range svcInfo.Edges {
			otherService := app.GetServiceByName(dep)
			service.AddDependency(otherService)
		}
	}

	// parse databases
	for _, dsInfo := range datastoresInfo {
		database := backends.NewDatabase(dsInfo.Name)
		app.AddDatabase(database)
	}

	// parse entrypoints
	for _, serviceName := range frontends {
		service := app.GetServiceByName(serviceName)
		for _, method := range service.GetMethods() {
			app.AddEntrypoint(service, method)
		}
	}
	for _, service := range app.GetAllServices() {
		if service.HasInitializerMethod() {
			// Run() method can also be considered as entrypoint
			// because they are always called when initializing services
			app.AddEntrypoint(service, "Run")
		}
	}
}

func (app *App) InitFields(pkgs []*ssa.Package) {
	for _, pkg := range pkgs {
		fmt.Printf("[APP PARSER] analyzing package: %s\n", pkg.String())
		for _, member := range pkg.Members {
			if ssaType, ok := member.(*ssa.Type); ok {
				service := app.GetServiceWithImplPathIfExists(ssaType.String())
				if service == nil {
					continue
				}
				fmt.Printf("\t[APP PARSER] found service impl: %s\n", service.GetImpl())
				if typeNamed, ok := ssaType.Type().(*types.Named); ok {
					if typeStruct, ok := typeNamed.Underlying().(*types.Struct); ok {
						i := 0
						for i < typeStruct.NumFields() {
							typeVar := typeStruct.Field(i)
							field := services.NewField(i, typeVar.Name())
							service.AddField(field)
							fmt.Printf("\t\t[APP PARSER] created new field: %s\n", field.String())
							i++
						}
					}
				}
			}
		}
	}
}
