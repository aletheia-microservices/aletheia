package constraints

import (
	"strings"

	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/auxten/postgresql-parser/pkg/walk"

	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
	//"analyzer/pkg/utils"
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

func parseSQLStatement(database *datastores.Datastore, sql string) {
	logger.Logger.Infof("[SQL PARSER] parsing statement: %s", sql)

	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	sql = strings.TrimSpace(sql)

	var tableName string
	var fields map[string]*datastores.Field

	w := &walk.AstWalker{
		Fn: func(ctx interface{}, node interface{}) (stop bool) {
			logger.Logger.Debugf("[SQL PARSER] visiting node (%T): %v", node, node)

			switch stmt := node.(type) {
			case *tree.CreateTable:
				tableName = stmt.Table.Table()
				fields = make(map[string]*datastores.Field, 0)
				
			case *tree.ColumnTableDef:
				columnName := stmt.Name.String()
				fieldName := tableName + "." + columnName
				fieldType := stmt.Type.Name()

				field := database.GetSchema().GetFieldIfExists(fieldName)
				if field == nil {
					field = datastores.NewField(fieldName, fieldType, -1, database)
					database.GetSchema().AddField(field)
					logger.Logger.Warnf("[SQL PARSER] added new database field: %s", field.GetFullName())
				}
				fields[columnName] = field

				for _, checkExpr := range stmt.CheckExprs {
					if comparisonExpr, ok := checkExpr.Expr.(*tree.ComparisonExpr); ok {
						constraint := datastores.NewConstraintNumerical(datastores.NewNumericalConstraint(
							comparisonExpr.Right.String(),
							datastores.ComparisonOperator(comparisonExpr.Operator),
						), field)
						field.AddConstraint(constraint)
						database.GetSchema().AddConstraint(constraint)
						logger.Logger.Warnf("[SQL PARSER] added new constraint %s", constraint.String())
					}
				}

			case *tree.UniqueConstraintTableDef:
				if stmt.PrimaryKey {
					constraint := datastores.NewConstraintPrimary()
					for _, column := range stmt.Columns {
						field := fields[column.Column.Normalize()]
						constraint.AddField(field)
						field.AddConstraint(constraint)
					}
					database.GetSchema().AddConstraint(constraint)
					logger.Logger.Warnf("[SQL PARSER] added new constraint: %s", constraint.String())
				}

			}
			return false
		},
	}
	
	stmts, err := parser.Parse(sql)
	if err != nil {
		logger.Logger.Fatalf("[SQL PARSER] %s", err.Error())
		return
	}
	
	ok, err := w.Walk(stmts, nil)
	if err != nil {
		logger.Logger.Fatalf("[SQL PARSER] %s", err.Error())
	} else if !ok{
		logger.Logger.Fatalf("[SQL PARSER] UNEXPECTED!")
	}

	//summarize(app, "SQL_PARSER")
}
