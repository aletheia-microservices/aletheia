package backends

import (
	"encoding/json"
)

type Database struct {
	name       string
	schemas    []*Schema
	typeString string // RelationalDB, Cache, NoSQLDatabase, Queue
}

func NewDatabase(name string, typeString string, schemas ...*Schema) *Database {
	return &Database{
		name:       name,
		schemas:    schemas,
		typeString: typeString,
	}
}

func (database *Database) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name   string    `json:"name"`
		Schema []*Schema `json:"schemas"`
	}{
		Name:   database.name,
		Schema: database.schemas,
	})
}

func (database *Database) GetName() string {
	return database.name
}

func (database *Database) AddSchema(schema *Schema) {
	database.schemas = append(database.schemas, schema)
}

func (database *Database) HasSchema(name string) bool {
	for _, schema := range database.schemas {
		if schema.name == name {
			return true
		}
	}
	return false
}

// name can be, for example, nosql collection, sql table, queue topic, etc.
func (database *Database) GetSchemaByNameIfExists(name string) *Schema {
	for _, schema := range database.schemas {
		if schema.GetName() == name {
			return schema
		}
	}
	return nil
}

func (database *Database) GetLastSchema() *Schema {
	return database.schemas[0]
}

func (database *Database) GetSchemas() []*Schema {
	return database.schemas
}

func (database *Database) IsQueue() bool {
	return database.typeString == "Queue"
}

func (database *Database) String() string {
	var str string
	str += " // schema: \n"
	for i, schema := range database.schemas {
		str += schema.String()
		if i < len(database.schemas)-1 {
			str += "\n"
		}
	}
	return database.name + str
}
