package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/app/services"
	"analyzer/pkg/utils"
)

type App struct {
	name        string
	databases   map[string]*backends.Database
	services    map[string]*services.Service
	entrypoints map[*services.Service][]string
}

func NewApp(name string) *App {
	return &App{
		name:        name,
		databases:   make(map[string]*backends.Database),
		services:    make(map[string]*services.Service),
		entrypoints: make(map[*services.Service][]string),
	}
}

func (app *App) GetAllFieldsReferencingCurrent(target *backends.Field) []*backends.Field {
	var fields []*backends.Field
	for _, database := range app.GetAllDatabases() {
		for _, schema := range database.GetAllSchemas() {
			for _, field := range schema.GetFields() {
				if field.HasConstraintForeignKeyToField(target) {
					fields = append(fields, field)
				}
			}
		}
	}
	return fields
}

// used by detectors
func (app *App) ComputeDatabaseFieldFromPath(database *backends.Database, fieldpath string) *backends.Field {
	schema := database.GetSchemaByNameIfExists(utils.ExtractSchemaNameFromFieldPath(fieldpath))
	if schema == nil {
		schemaName := utils.ExtractSchemaNameFromFieldPath(fieldpath)
		if strings.HasSuffix(schemaName, "[*]") {
			// [TO BE IMPROVED]
			// sometimes we get schema name "schema[*]" from fieldpaths "schema[*].Value"
			// because mongodb read filter fields are not being yet parsed for reads taints
			// for now we hardcode to remove the [*] in "schema[*]"
			schemaName = schemaName[:len(schemaName)-3]
			schema = database.GetSchemaByNameIfExists(schemaName)
		} else {
			log.Panicf("[APP] nil schema (%s) for fieldpath (%s)\n", utils.ExtractSchemaNameFromFieldPath(fieldpath), fieldpath)
		}
	}
	//EVAL - fmt.Printf("[APP] get field for (%s) in schema (%s)\n", fieldpath, schema.GetName())
	// [TO BE IMPROVED]
	field := schema.GetFieldByPathIfExists(fieldpath)

	// [TO BE IMPROVED]
	// in the future, the ssa parser should be the one to infer
	// all schema fields (from AST structure) beforehand
	if field == nil {
		field = backends.NewField(fieldpath, database, schema)
		schema.AddField(field)
	}
	return field
}

func (app *App) GetAllDatabaseFieldsWithPrefixPath(field *backends.Field, include bool) []*backends.Field {
	var matchingFields []*backends.Field
	schema := field.GetSchema()
	for _, otherField := range schema.GetAllFieldsLst() {
		if otherField == field && !include {
			continue
		}
		if strings.HasPrefix(otherField.GetPath(), field.GetPath()) {
			matchingFields = append(matchingFields, otherField)
		}
	}
	return matchingFields
}

func (app *App) GetName() string {
	return app.name
}

func (app *App) SetServiceEntrypoints(service *services.Service, methods []string) {
	sort.Strings(methods)
	app.entrypoints[service] = methods
}

func (app *App) AddEntrypoint(service *services.Service, method string) {
	app.entrypoints[service] = append(app.entrypoints[service], method)
}

// e.g., postnotification.UploadService.UploadPost
func (app *App) GetEntrypointsShortPaths() []string {
	//EVAL - fmt.Printf("[APP] getting entrypoint short paths...\n")
	var entrypoints []string
	for service, serviceEntrypoints := range app.entrypoints {
		for _, method := range serviceEntrypoints {
			//EVAL - fmt.Printf("\t[APP] found (%s)\n", service.GetPackage()+"."+service.GetName()+"."+method)
			entrypoints = append(entrypoints, service.GetPackage()+"."+service.GetName()+"."+method)
		}
	}
	return entrypoints
}

func (app *App) GetServiceWithPathIfExists(path string) *services.Service {
	for _, service := range app.services {
		if service.GetPath() == path {
			return service
		}
	}
	return nil
}

func (app *App) GetServiceWithImplPath(implpath string) *services.Service {
	for _, service := range app.services {
		if service.GetPackagePath()+"."+service.GetImpl() == implpath {
			return service
		}
	}
	/* for _, service := range app.services {
		//EVAL - fmt.Printf("pkg path = (%s) // impl = (%s)\n", service.GetPackagePath(), service.GetImpl())
	} */
	log.Fatalf("could not find service for impl path (%s)", implpath)
	return nil
}

func (app *App) GetServiceWithImplPathIfExists(implpath string) *services.Service {
	for _, service := range app.services {
		if service.GetPackagePath()+"."+service.GetImpl() == implpath {
			return service
		}
	}
	return nil
}

func (app *App) GetServiceWithConstructorShortPathIfExists(constructorPath string) *services.Service {
	for _, service := range app.services {
		if service.GetPackage()+"."+service.GetConstructor() == constructorPath {
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

func (app *App) GetAllDatabases() []*backends.Database {
	var lst []*backends.Database
	for _, db := range app.databases {
		lst = append(lst, db)
	}
	return lst
}

func (app *App) GetDatabaseByName(name string) *backends.Database {
	db, ok := app.databases[name]
	if !ok {
		log.Fatalf("could not find database for name (%s) in app", name)
	}
	return db
}

func (app *App) GetAllServices() []*services.Service {
	var lst []*services.Service
	for _, service := range app.services {
		lst = append(lst, service)
	}
	return lst
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
	// sort databases
	sort.Strings(databasesLst)

	servicesLst := make([]*services.Service, len(app.services))
	i = 0
	for _, service := range app.services {
		servicesLst[i] = service
		i++
	}

	// sort services
	sort.Slice(servicesLst, func(i, j int) bool {
		return servicesLst[i].GetName() < servicesLst[j].GetName()
	})

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
	var sortedDbs []*backends.Database
	for _, db := range app.databases {
		sortedDbs = append(sortedDbs, db)
	}
	sort.Slice(sortedDbs, func(i, j int) bool {
		return sortedDbs[i].GetName() < sortedDbs[j].GetName()
	})

	for _, db := range sortedDbs {
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
