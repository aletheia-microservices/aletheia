// Package backends represents databases and their entities (schemas), fields and constraints
package backends

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
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
	sortedSchemas := database.schemas
	sort.Slice(sortedSchemas, func(i, j int) bool {
		return sortedSchemas[i].GetName() < sortedSchemas[j].GetName()
	})
	return json.Marshal(&struct {
		Name   string    `json:"name"`
		Schema []*Schema `json:"schemas"`
	}{
		Name:   database.name,
		Schema: sortedSchemas,
	})
}

func (database *Database) GetAllSchemas() []*Schema {
	return database.schemas
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

func (database *Database) GetOrCreateSchema(name string) *Schema {
	for _, schema := range database.schemas {
		if schema.GetName() == name {
			return schema
		}
	}
	schema := NewSchema(name, database)
	database.AddSchema(schema)
	return schema
}

// GetSchemaForFieldPath returns the schema of a field path (<database>.<schema>[.<field>]),
// e.g., user_db.address.City => address
func (database *Database) GetSchemaForFieldPath(fieldpath string) *Schema {
	parts := strings.SplitN(fieldpath, ".", 3)
	if len(parts) < 2 {
		return nil
	}
	schemaName := parts[1]
	// TODO(improvement): sometimes we get schema name "schema[*]" from fieldpaths "schema[*].Value"
	// because mongodb read filter fields are not being yet parsed for reads taints
	// for now we hardcode to remove the [*] in "schema[*]"
	schemaName = strings.TrimSuffix(schemaName, "[*]")
	schema := database.GetSchemaByNameIfExists(schemaName)
	if schema == nil {
		// the schema of every database call is created when parsing it, so this means the field path
		// does not name a schema of the database (e.g., travel_db.tripResponse.PriceForConfortClass in trainticket,
		// where tripResponse is a field of the struct that holds the trip)
		logrus.Warnf("[DATABASE] creating schema (%s) for field path (%s) not matching any schema of database (%s)", schemaName, fieldpath, database.name)
		schema = database.GetOrCreateSchema(schemaName)
	}
	return schema
}

func (database *Database) GetSchemas() []*Schema {
	return database.schemas
}

func (database *Database) IsQueue() bool {
	return database.typeString == "Queue"
}

func (database *Database) IsNoSQL() bool {
	return database.typeString == "NoSQLDatabase"
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
