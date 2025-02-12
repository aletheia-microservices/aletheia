package constraints

import (
	"strings"

	"github.com/xwb1989/sqlparser"

	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

type SQLTable struct {
	Name        string
	Columns     []SQLColumn
	PrimaryKeys []string
}
type SQLColumn struct {
	Name         string
	Type         string
	DefaultValue string
	IsPrimaryKey bool
}

func parseSQLStatement(app *app.App, database *datastores.Datastore, sql string) {
	logger.Logger.Infof("[SQL PARSER] parsing statement: %s", sql)

	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	sql = strings.TrimSpace(sql)

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		panic(err)
	}

	switch stmt := stmt.(type) {
	case *sqlparser.DDL:
		logger.Logger.Tracef("[SQL PARSER] parsing DDL: %v", stmt)
		tableName := stmt.NewName.Name.CompliantName()
		if stmt.TableSpec == nil {
			logger.Logger.Fatalf("[SQL PARSER] nil tablespec for SQL statament: %v", stmt)
		}
		fields := make(map[string]*datastores.Field, 0)
		for _, column := range stmt.TableSpec.Columns {
			fieldName := column.Name.CompliantName()

			columnName := tableName + "." + fieldName
			columnType := column.Type.Type
			field := datastores.NewField(columnName, columnType, -1, database)
			database.GetSchema().AddField(field)

			fields[column.Name.CompliantName()] = field
			logger.Logger.Infof("[SQL PARSER] added new database field: %s", field.GetFullName())
		}
		for _, index := range stmt.TableSpec.Indexes {
			var constraint *datastores.Constraint
			if index.Info.Unique {
				constraint = datastores.NewConstraintUnique()
			} else if index.Info.Primary {
				constraint = datastores.NewConstraintPrimary()
			}
			for _, column := range index.Columns {
				field := fields[column.Column.CompliantName()]
				constraint.AddField(field)
				field.AddConstraint(constraint)

			}
			database.GetSchema().AddConstraint(constraint)
			logger.Logger.Warnf("[SQL PARSER] added new constraint: %s", constraint.String())
		}
	default:
		logger.Logger.Fatalf("[SQL PARSER] unexpected type (%s) for sqlparser: %v", utils.GetType(stmt), stmt)
	}
	summarize(app, "SQL_PARSER")
}
