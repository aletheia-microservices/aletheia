package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/app/services"
)

type App struct {
	name      string
	databases map[string]*backends.Database
	services  map[string]*services.Service
}

func NewApp(name string) *App {
	return &App{
		name:      name,
		databases: make(map[string]*backends.Database),
		services:  make(map[string]*services.Service),
	}
}

func (app *App) GetName() string {
	return app.name
}

func (app *App) GetServiceWithPathIfExists(path string) *services.Service {
	for _, service := range app.services {
		if service.GetPath() == path {
			return service
		}
	}
	return nil
}

func (app *App) AddDatabase(db *backends.Database) {
	app.databases[db.GetName()] = db
}

func (app *App) HasDatabase(name string) bool {
	_, ok := app.databases[name]
	return ok
}

func (app *App) GetDatabaseByName(name string) *backends.Database {
	db, ok := app.databases[name]
	if !ok {
		log.Fatalf("could not find database for name (%s) in app", name)
	}
	return db
}

func (app *App) AddService(service *services.Service) {
	app.services[service.GetName()] = service
}

func (app *App) HasService(name string) bool {
	_, ok := app.services[name]
	return ok
}

func (app *App) GetServiceByName(name string) *services.Service {
	service, ok := app.services[name]
	if !ok {
		log.Fatalf("could not find service for name (%s) in app", name)
	}
	return service
}

func (app *App) String() string {
	str := "APP (" + app.name + ")\n"
	str += "\n=== SERVICES ===\n"
	i := 0
	for _, service := range app.services {
		str += service.String()
		if i < len(app.services)-1 {
			str += "\n"
		}
		i++
	}
	str += "\n\n=== DATABASES ===\n"
	i = 0
	for _, db := range app.databases {
		str += db.String()
		if i < len(app.databases)-1 {
			str += "\n"
		}
		i++
	}
	return str
}

func (app *App) MarshalJSON() ([]byte, error) {
	databasesLst := make([]string, len(app.databases))
	i := 0
	for _, db := range app.databases {
		databasesLst[i] = db.GetName()
		i++
	}
	servicesLst := make([]*services.Service, len(app.services))
	i = 0
	for _, service := range app.services {
		servicesLst[i] = service
		i++
	}
	return json.Marshal(&struct {
		Name      string              `json:"name"`
		Databases []string            `json:"databases"`
		Services  []*services.Service `json:"services"`
	}{
		Name:      app.name,
		Databases: databasesLst,
		Services:  servicesLst,
	})
}

func (app *App) WriteAppToJSON() error {
	filename := fmt.Sprintf("output/%s/app.json", app.name)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (app *App) WriteSchemaToJSON() error {
	filename := fmt.Sprintf("output/%s/schema.json", app.name)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	dbJSONMap := make(map[string]json.RawMessage, len(app.databases))
	for _, db := range app.databases {
		data, err := json.Marshal(db)
		if err != nil {
			return fmt.Errorf("failed to marshal database %s: %w", db.GetName(), err)
		}
		dbJSONMap[db.GetName()] = data
	}

	data, err := json.MarshalIndent(dbJSONMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal combined schema: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}
