package registry

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
)

// extract the bson name from a struct field tag string
// examples:
// - bson:"_id" 		-> returns ("_id", true)
// - bson:"Title" 		-> returns ("Title", true).
// - untagged fields 	-> default to strings.ToLower(fieldName).
func parseBsonTagName(tag string) (string, bool) {
	const prefix = `bson:"`
	idx := strings.Index(tag, prefix)
	if idx == -1 {
		return "", false
	}
	rest := tag[idx+len(prefix):]
	end := strings.Index(rest, `"`)
	if end == -1 {
		return "", false
	}
	name := rest[:end]
	// remove things like ",omitempty"
	if comma := strings.Index(name, ","); comma != -1 {
		name = name[:comma]
	}
	return name, true
}

func RegisterNoSQLPrimaryKey(app *app.App, databaseStr string, collectionStr string, docVal ssa.Value) {
	database := app.GetDatabaseByName(databaseStr)
	schema := database.GetOrCreateSchema(collectionStr)
	if makeIface, ok := docVal.(*ssa.MakeInterface); ok {
		if named, ok := makeIface.X.Type().(*types.Named); ok {
			if str, ok := named.Underlying().(*types.Struct); ok {
				for i := range str.NumFields() {
					fieldVar := str.Field(i)
					tag := str.Tag(i)

					bsonName, hasBson := parseBsonTagName(tag)
					if !hasBson {
						// mongodb driver lowercases untagged fields
						bsonName = strings.ToLower(fieldVar.Name())
					}

					if hasBson && bsonName == "_id" {
						bsonFieldpath := database.GetName() + "." + schema.GetName() + "." + bsonName
						goFieldpath := database.GetName() + "." + schema.GetName() + "." + fieldVar.Name()
						field := schema.GetOrCreateField(database, bsonFieldpath)
						schema.AddBsonFieldAlias(goFieldpath, bsonFieldpath)
						constraint := backends.NewConstraint(backends.CONSTRAINT_PRIMARY, field)
						field.AddConstraint(constraint)
						schema.AddConstraint(constraint)
					}
				}
			}
		}
	}
}
