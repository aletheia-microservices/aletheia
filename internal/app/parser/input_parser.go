package appparser

import (
	"slices"

	abstractgraphinput "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/input"
	"github.com/aletheia-microservices/aletheia/internal/app"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
	"github.com/aletheia-microservices/aletheia/internal/app/services"
)

// InitFromInput is the equivalent of Init (and of the schema parsers) for an app described by an
// input model instead of Blueprint: it adds the databases with their keys and the services of model,
// which must be valid (see abstractgraphinput.LoadInputModel)
func InitFromInput(app *app.App, model *abstractgraphinput.InputModel) {
	// parse databases
	for _, db := range model.Databases {
		database := backends.NewDatabase(db.Name, db.Type)
		app.AddDatabase(database)
		for _, s := range db.Schemas {
			schema := database.GetOrCreateSchema(s.Name)
			if len(s.PrimaryKey) > 0 {
				constraint := backends.NewConstraint(backends.CONSTRAINT_PRIMARY)
				for _, name := range s.PrimaryKey {
					field := schema.GetOrCreateField(database, db.Name+"."+s.Name+"."+name)
					constraint.AddField(field)
					field.AddConstraint(constraint)
				}
				schema.AddConstraint(constraint)
			}
			for _, name := range s.Unique {
				field := schema.GetOrCreateField(database, db.Name+"."+s.Name+"."+name)
				constraint := backends.NewConstraint(backends.CONSTRAINT_UNIQUE, field)
				field.AddConstraint(constraint)
				schema.AddConstraint(constraint)
			}
		}
	}

	// parse services
	for _, fn := range model.Functions {
		if !app.HasService(fn.Service) {
			app.AddService(services.NewService(fn.Service, "", "", "", "", "", nil))
		}
		service := app.GetServiceByName(fn.Service)
		if !service.HasMethod(fn.Method) {
			service.SetMethods(fn.Method)
		}
	}
	deps := make(map[*services.Service][]*services.Service)
	for _, fn := range model.Functions {
		service := app.GetServiceByName(fn.Service)
		for _, call := range fn.Calls {
			if call.ServiceCall == nil {
				continue
			}
			dep := app.GetServiceByName(call.ServiceCall.Service)
			if dep != service && !slices.Contains(deps[service], dep) {
				deps[service] = append(deps[service], dep)
				service.AddDependency(dep)
			}
		}
	}
}
