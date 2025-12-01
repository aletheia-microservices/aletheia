package registry

import (
	"go/types"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
)

func RegisterNoSQLPrimaryKey(app *app.App, databaseStr string, collectionStr string, docVal ssa.Value) {
	database := app.GetDatabaseByName(databaseStr)
	schema := database.GetOrCreateSchema(collectionStr)
	if makeIface, ok := docVal.(*ssa.MakeInterface); ok {
		if named, ok := makeIface.X.Type().(*types.Named); ok {
			if str, ok := named.Underlying().(*types.Struct); ok {
				for i := range str.NumFields() {
					fieldVar := str.Field(i)
					tag := str.Tag(i)
					if tag == "bson:\"_id\"" {
						fieldpath := database.GetName() + "." + schema.GetName() + "." + fieldVar.Name()
						field := schema.GetOrCreateField(database, fieldpath)
						constraint := backends.NewConstraint(backends.CONSTRAINT_PRIMARY, field)
						field.AddConstraint(constraint)
						schema.AddConstraint(constraint)
						logrus.Tracef("[CALLS BLUEPRINT] [NOSQL PK] registered primary key constraint: %s\n", constraint.String())
					}
				}
			}
		}
	}
	//logrus.Tracef("[CALLS BLUEPRINT] [NOSQL PK] skipping registerNoSQLPrimaryKey (database=%s, schema=%s)\n", database.GetName(), schema.GetName())
}
