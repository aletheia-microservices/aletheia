package constraints

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
)

type Property struct {
	BsonType    string `json:"bsonType"`
	Description string `json:"description"`
	Minimum     *int   `json:"minimum,omitempty"`
	Maximum     *int   `json:"maximum,omitempty"`
	MaxLength   *int   `json:"maxLength,omitempty"`
}

type Schema struct {
	Properties map[string]Property `json:"properties"`
	Database   string
	Collection string
}

type JSONSchema struct {
	Schema Schema `json:"$jsonSchema"`
}

// Parse SQL files and return slice of SQL statements
func parseAppDatabaseDocSchemas(input string) []JSONSchema {
	var jsonSchemas []JSONSchema
	targetDbPaths := strings.Split(input, ";")
	for _, dbPath := range targetDbPaths {
		splits := strings.Split(dbPath, ":")
		databaseName, collectionName, jsonFilePath := splits[0], splits[1], splits[2]
		jsonBytes, err := os.ReadFile(jsonFilePath)
		if err != nil {
			logger.Logger.Fatalf("error reading json files: %s", err.Error())
			return nil
		}

		var jsonSchema JSONSchema
		err = json.Unmarshal([]byte(jsonBytes), &jsonSchema)
		if err != nil {
			logger.Logger.Fatalf("[DOC PARSER] error parsing json data: %s", err.Error())
		}

		jsonSchema.Schema.Database = databaseName
		jsonSchema.Schema.Collection = collectionName
		jsonSchemas = append(jsonSchemas, jsonSchema)
	}
	return jsonSchemas
}

func parseDocJSON(database *datastores.Datastore, jsonSchema JSONSchema) {
	logger.Logger.Infof("[DOC PARSER] parsing json data: %v", jsonSchema)

	for fieldName, prop := range jsonSchema.Schema.Properties {
		fieldFullName := jsonSchema.Schema.Collection + "." + fieldName
		field := database.GetSchema().GetFieldIfExists(fieldFullName)
		if field == nil {
			field = datastores.NewField(fieldFullName, prop.BsonType, -1, database)
			database.GetSchema().AddField(field)
			logger.Logger.Warnf("[DOC PARSER] added new database field: %s", field.GetFullName())
		} else {
			logger.Logger.Debugf("[DOC PARSER] got field: %v", field.String())
		}

		var numericalValue string
		var numericalComparatorOp datastores.ComparisonOperator
		if prop.Minimum != nil {
			numericalComparatorOp = datastores.GE
			numericalValue = fmt.Sprintf("%d", *prop.Minimum)
		}
		if prop.Maximum != nil {
			numericalComparatorOp = datastores.LE
			numericalValue = fmt.Sprintf("%d", *prop.Maximum)
		}

		if numericalValue != "" {
			constraint := datastores.NewConstraintNumerical(datastores.NewNumericalConstraint(numericalValue, numericalComparatorOp), field)
			field.AddConstraint(constraint)
			database.GetSchema().AddConstraint(constraint)
			logger.Logger.Warnf("[DOC PARSER] added new constraint %s", constraint.String())
		}
	}
}
