package app

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/config"
	"analyzer/pkg/utils"
)

type Property struct {
	BsonType    string `json:"bsonType"`
	Description string `json:"description"`
	Minimum     *int   `json:"minimum,omitempty"`
	Maximum     *int   `json:"maximum,omitempty"`
	MaxLength   *int   `json:"maxLength,omitempty"`
}

type Schema struct {
	Properties  map[string]Property `json:"properties"`
	UniqueItems []string            `json:"uniqueItems"`
	Database    string              `json:"database"`
	Collection  string              `json:"collection"`
}

type JSONSchema struct {
	Schema Schema `json:"$jsonSchema"`
}

func (app *App) ParseNoSQLSchemaFromUserFile() {
	pkgPathDB := utils.APPS_NOSQL_SCHEMAS[app.name] + "/database/"
	filepaths, err := filepath.Glob(filepath.Join(pkgPathDB, "*.json"))
	if err != nil {
		logrus.Fatalf("error extracting database json files: %s", err.Error())
	}

	var jsonSchemas []JSONSchema
	for _, jsonFilePath := range filepaths {
		jsonBytes, err := os.ReadFile(jsonFilePath)
		if err != nil {
			logrus.Fatalf("error reading json files: %s", err.Error())
		}

		var jsonSchema JSONSchema
		err = json.Unmarshal([]byte(jsonBytes), &jsonSchema)
		if err != nil {
			logrus.Fatalf("[JSON PARSER] error parsing json data: %s", err.Error())
		}

		jsonSchemas = append(jsonSchemas, jsonSchema)
	}

	for _, jsonSchema := range jsonSchemas {
		parseJson(app, jsonSchema)
	}
}

func parseJson(app *App, jsonSchema JSONSchema) {
	for _, field := range jsonSchema.Schema.UniqueItems {
		database := app.GetDatabaseByName(jsonSchema.Schema.Database)
		schema := database.GetOrCreateSchema(jsonSchema.Schema.Collection)
		fieldpath := jsonSchema.Schema.Database + "." + jsonSchema.Schema.Collection + "." + field
		field := schema.GetOrCreateField(database, fieldpath)
		constraint := backends.NewConstraint(backends.CONSTRAINT_UNIQUE, field)
		field.AddConstraint(constraint)
		schema.AddConstraint(constraint)

		if config.Global.MakeIndexesAsPrimaryKeysForNoSQLDatabases {
			pkConstraint := backends.NewConstraint(backends.CONSTRAINT_PRIMARY, field)
			field.AddConstraint(pkConstraint)
			schema.AddConstraint(pkConstraint)
		}
	}
}
