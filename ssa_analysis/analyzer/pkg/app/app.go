package app

import (
	"log"

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
		if i < len(app.databases) - 1 {
			str += "\n"
		}
		i++
	}
	return str
}
