package constraints

import (
	"strings"

	"github.com/xwb1989/sqlparser"

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

func ParseSQLStatement(database *datastores.Datastore, sql string) {
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
		logger.Logger.Debugf("[SQL PARSER] parsing DDL: %v", stmt)
		tableName := stmt.NewName.Name.CompliantName()
		if stmt.TableSpec == nil {
			logger.Logger.Fatalf("[SQL PARSER] nil tablespec for SQL statament: %v", stmt)
		}
		fields := make(map[string]*datastores.Field, 0)
		for _, column := range stmt.TableSpec.Columns {
			fieldName := column.Name.CompliantName()

			columnName := tableName + "." + fieldName
			columnType := column.Type.Type
			field := datastores.NewEntry(columnName, columnType, -1, database)
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

				logger.Logger.Warnf("[SQL PARSER] added new constraint: %s", constraint.String())
			}
		}
	default:
		logger.Logger.Fatalf("[SQL PARSER] unexpected type (%s) for sqlparser: %v", utils.GetType(stmt), stmt)
	}
}

/* func ParseSQLStatement2(database *Datastore, stmt string) {
	logger.Logger.Infof("parsing sql statement @ %s: %s", database.GetName(), stmt)

	stmt = strings.ReplaceAll(stmt, "\n", " ")
	stmt = strings.ReplaceAll(stmt, "\t", " ")
	stmt = strings.TrimSpace(stmt)

	tableRegex := regexp.MustCompile(`(?i)CREATE TABLE IF NOT EXISTS (\w+)`)
	tableMatch := tableRegex.FindStringSubmatch(stmt)
	if len(tableMatch) < 2 {
		logger.Logger.Fatalf("[SQL] failed to extract table name for SQL statement: %s", stmt)
	}
	tableName := tableMatch[1]

	columnRegex := regexp.MustCompile(`$begin:math:text$(.*)$end:math:text$\s*$`)
	columnMatch := columnRegex.FindStringSubmatch(stmt)
	if len(columnMatch) < 2 {
		logger.Logger.Fatalf("[SQL] failed to extract column definitions: %s", stmt)
	}

	columnsRaw := columnMatch[1]
	columnLines := strings.Split(columnsRaw, ",")

	var columns []SQLColumn
	var primaryKeys []string

	for _, line := range columnLines {
		line = strings.TrimSpace(line)

		// check if primary key
		if strings.HasPrefix(strings.ToUpper(line), "PRIMARY KEY") {
			pkRegex := regexp.MustCompile(`(?i)PRIMARY KEY\s*$begin:math:text$(.*)$end:math:text$`)
			pkMatch := pkRegex.FindStringSubmatch(line)
			if len(pkMatch) > 1 {
				primaryKeys = strings.Split(pkMatch[1], ",")
				for i := range primaryKeys {
					primaryKeys[i] = strings.TrimSpace(primaryKeys[i])
				}
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		var defaultVal string
		for i, part := range parts {
			if strings.ToUpper(part) == "DEFAULT" && i+1 < len(parts) {
				defaultVal = parts[i+1]
			}
		}

		columns = append(columns, SQLColumn{
			Name:         parts[0],
			Type:         parts[1],
			DefaultValue: defaultVal,
			IsPrimaryKey: false,
		})
	}

	for i := range columns {
		for _, pk := range primaryKeys {
			if columns[i].Name == pk {
				columns[i].IsPrimaryKey = true
			}
		}
	}

	table := SQLTable{
		Name:        tableName,
		Columns:     columns,
		PrimaryKeys: primaryKeys,
	}

	logger.Logger.Infof("GOT TABLE: %v", table)
} */
