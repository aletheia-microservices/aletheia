package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"analyzer/pkg/app/backends"
)

type App struct {
	name      string
	databases map[string]*backends.Database
}

func NewApp(name string) *App {
	return &App{
		name:      name,
		databases: make(map[string]*backends.Database),
	}
}

func (app *App) GetName() string {
	return app.name
}

func (app *App) MarshalJSON() ([]byte, error) {
	databases := make([]*backends.Database, len(app.databases))
	i := 0
	for _, db := range app.databases {
		databases[i] = db
		i++
	}
	return json.Marshal(&struct {
		Name      string               `json:"name"`
		Databases []*backends.Database `json:"databases"`
	}{
		Name:      app.name,
		Databases: databases,
	})
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

func (app *App) String() string {
	str := "APP (" + app.name + ")\n"
	str += "=== DATABASES ===\n"
	i := 0
	for _, db := range app.databases {
		str += db.String()
		if i < len(app.databases)-1 {
			str += "\n"
		}
		i++
	}
	return str
}

func (app *App) WriteToJSON() error {
	filename := fmt.Sprintf("output/%s/data-schema.json", app.name)
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
