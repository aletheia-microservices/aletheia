package app

import (
	"os"
	"strings"

	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/auxten/postgresql-parser/pkg/walk"
	"github.com/sirupsen/logrus"
	"github.com/xwb1989/sqlparser"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/utils"
)

// -----------
// FILE PARSER
// -----------

type sqlTable struct {
	Name        string
	Columns     []sqlColumn
	PrimaryKeys []string
}
type sqlColumn struct {
	Name         string
	Type         string
	DefaultValue string
	IsPrimaryKey bool
}

type sqlDbStmt struct {
	db   string
	stmt string
}

func (app *App) ParseSQLSchemaFromUserFile() {
	if ok, input := utils.GetAppDatabaseSQLPaths(app.GetName(), true); ok {
		dbStmts := parseSQLStatementsFromInput(input)
		for _, dbStmt := range dbStmts {
			database := app.GetDatabaseByName(dbStmt.db)
			parseSQLStatement(database, dbStmt.stmt)
		}
	}
}

// Parse SQL files and return slice of SQL statements
func parseSQLStatementsFromInput(input string) []sqlDbStmt {
	var dbStmts []sqlDbStmt
	targetDbPaths := strings.Split(input, ";")
	for _, dbPath := range targetDbPaths {
		splits := strings.Split(dbPath, ":")
		db := splits[0]
		sqlStmt := splits[1]
		sqlBytes, err := os.ReadFile(sqlStmt)
		if err != nil {
			logrus.Fatalf("error reading sql files: %s", err.Error())
			return nil
		}
		sqlStmts := strings.Split(string(sqlBytes), ";")
		for _, stmt := range sqlStmts {
			if stmt == "\n" {
				continue
			}
			dbStmts = append(dbStmts, sqlDbStmt{db, stmt})
		}
	}
	return dbStmts
}

func parseSQLStatement(database *backends.Database, sql string) {
	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	sql = strings.TrimSpace(sql)

	var tableName string
	var fields map[string]*backends.Field

	w := &walk.AstWalker{
		Fn: func(ctx interface{}, node interface{}) (stop bool) {
			switch stmt := node.(type) {
			case *tree.CreateTable:
				tableName = stmt.Table.Table()
				fields = make(map[string]*backends.Field, 0)

			case *tree.ColumnTableDef:
				columnName := stmt.Name.String()
				fieldName := tableName + "." + columnName
				fieldPath := database.GetName() + "." + fieldName

				schema := database.GetOrCreateSchema(tableName)

				field := schema.GetOrCreateField(database, fieldPath)
				fields[columnName] = field

				var constraint *backends.Constraint
				if stmt.PrimaryKey.IsPrimaryKey {
					constraint = backends.NewConstraint(backends.CONSTRAINT_PRIMARY, field)
				} else if stmt.Unique {
					constraint = backends.NewConstraint(backends.CONSTRAINT_UNIQUE, field)
				}
				if constraint != nil {
					field.AddConstraint(constraint)
					schema.AddConstraint(constraint)
				}

			case *tree.UniqueConstraintTableDef:
				var schema *backends.Schema

				var constraint *backends.Constraint
				if stmt.PrimaryKey {
					constraint = backends.NewConstraint(backends.CONSTRAINT_PRIMARY)
				} else {
					constraint = backends.NewConstraint(backends.CONSTRAINT_UNIQUE)
				}

				for _, column := range stmt.Columns {
					field := fields[column.Column.Normalize()]
					constraint.AddField(field)
					field.AddConstraint(constraint)

					if schema == nil {
						schema = field.GetSchema()
					}
				}

				schema.AddConstraint(constraint)

			}
			return false
		},
	}

	stmts, err := parser.Parse(sql)
	if err != nil {
		logrus.Fatalf("[APP SQL PARSER] %s", err.Error())
		return
	}

	ok, err := w.Walk(stmts, nil)
	if err != nil {
		logrus.Fatalf("[APP SQL PARSER] %s", err.Error())
	} else if !ok {
		logrus.Fatalf("[APP SQL PARSER] UNEXPECTED!")
	}
}

// -----------
// CODE PARSER
// -----------

type tableNameAlias struct {
	alias string
	name  string
}

// ParseSQLRead returns two slices, one for written fields and another for filter fields
// where each field return has format <table>.<field>
//
// NOTE: for SQL Selects on all fields (i.e., '*') the readFields length is 1
// and the readField has format <database>.<table>
func ParseSQLRead(db string, stmtStr string) ([]string, []string, []string, bool) {
	stmt, err := sqlparser.Parse(stmtStr)
	if err != nil {
		logrus.Warnf("[APP SQL PARSER] unable to parse sql query (%s): %s", stmtStr, err.Error())
		return nil, nil, nil, false
	}

	var selectedFields []string
	var filterFields []string
	var tableNameAliasLst []tableNameAlias

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		tableNameAliasLst = parseSQLTableExprs(stmt.From)

		for _, expr := range stmt.SelectExprs {
			if _, ok := expr.(*sqlparser.StarExpr); ok {
				fieldpath := db + "." + tableNameAliasLst[0].name // + ".*"
				selectedFields = append(selectedFields, fieldpath)
			} else if aliasedExpr, ok := expr.(*sqlparser.AliasedExpr); ok {
				if valTuple, ok := aliasedExpr.Expr.(sqlparser.ValTuple); ok {
					for _, expr := range valTuple {
						if col, ok := expr.(*sqlparser.ColName); ok {
							prefixTableName, columnName := parseColumnName(string(col.Name.CompliantName()))
							tableName := parseTableName(prefixTableName, tableNameAliasLst)
							fieldpath := db + "." + tableName + "." + columnName

							selectedFields = append(selectedFields, fieldpath)
						}
					}
				}
			}
		}

		if stmt.Where != nil {
			var ok bool
			filterFields, ok = parseSQLWhere(db, stmt.Where, tableNameAliasLst)
			if !ok {
				logrus.Warnf("[APP SQL PARSER] unable to parse sql WHERE clause: %v", stmt.Where)
				return nil, nil, nil, false
			}
		}

	default:
		logrus.Fatalf("[APP SQL PARSER] Unsupported SQL statement: %s", stmtStr)
	}

	tableNames := make([]string, len(tableNameAliasLst))
	for i, tableAlias := range tableNameAliasLst {
		tableNames[i] = tableAlias.name
	}

	return selectedFields, filterFields, tableNames, true
}

func parseSQLWhere(db string, stmtWhere *sqlparser.Where, tableNameAliasLst []tableNameAlias) ([]string, bool) {
	var filterFields []string
	if comparisonExpr, ok := stmtWhere.Expr.(*sqlparser.ComparisonExpr); ok {
		var leftFieldName string
		if col, ok := comparisonExpr.Left.(*sqlparser.ColName); ok {
			prefixTableName, columnName := parseColumnName(string(col.Name.CompliantName()))
			tableName := parseTableName(prefixTableName, tableNameAliasLst)
			leftFieldName = db + "." + tableName + "." + columnName
		}
		filterFields = append(filterFields, leftFieldName)
	}
	return filterFields, true
}

// ParseSQLWrite returns two slices, one for written fields and another for filter fields
// where each field return has format <table>.<field>
func ParseSQLWrite(db string, stmtStr string) ([]string, []string, string, bool) {
	stmt, err := sqlparser.Parse(stmtStr)
	if err != nil {
		logrus.Fatalf("[APP SQL PARSER] unable to parse sql query (%s): %s", stmtStr, err.Error())
		return nil, nil, "", false
	}

	var writtenFields []string
	var filterFields []string
	var tableName string

	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		if values, ok := stmt.Rows.(sqlparser.Values); ok {
			for _, tuple := range values {
				for colIdx, expr := range tuple {
					if sqlVal, ok := expr.(*sqlparser.SQLVal); ok {
						tableName = stmt.Table.Name.CompliantName()
						if sqlVal.Type == sqlparser.ValArg { // placeholder (e.g., '?' that is then parsed into ':v1', ':v2', etc.)
							fieldpath := db + "." + tableName + "." + stmt.Columns[colIdx].CompliantName()
							writtenFields = append(writtenFields, fieldpath)
						}
					}
				}
			}
		} else {
			logrus.Fatalf("[APP SQL PARSER] unexpected type %T for rows in sql insert: %s", stmt.Rows, stmtStr)
		}
	case *sqlparser.Update:
		argIdx := 0
		tableNameAliasLst := parseSQLTableExprs(stmt.TableExprs)

		for _, expr := range stmt.Exprs {

			prefixTableName, columnName := parseColumnName(string(expr.Name.Name.CompliantName()))
			tableName = parseTableName(prefixTableName, tableNameAliasLst)
			fieldpath := db + "." + tableName + "." + columnName

			if sqlVal, ok := expr.Expr.(*sqlparser.SQLVal); ok {
				if sqlVal.Type == sqlparser.ValArg { // placeholder (e.g., '?', parsed as ':v1', ':v2')
					//fieldObj := args[argIdx]
					writtenFields = append(writtenFields, fieldpath)
					argIdx++
				}
			}
		}
		var ok bool
		filterFields, ok = parseSQLWhere(db, stmt.Where, tableNameAliasLst)
		if !ok {
			return nil, nil, "", false
		}
	}
	return writtenFields, filterFields, tableName, true
}

// ParseSQLDelete parses a DELETE statement and returns:
// - filterFields: fields used in the WHERE clause, in the form <db>.<table>.<column>
// - tableNames: the table names involved in the delete
func ParseSQLDelete(db string, stmtStr string) ([]string, []string, bool) {
	stmt, err := sqlparser.Parse(stmtStr)
	if err != nil {
		logrus.Fatalf("[APP SQL PARSER] unable to parse sql query (%s): %s", stmtStr, err.Error())
	}

	var filterFields []string
	var tableNameAliasLst []tableNameAlias

	switch stmt := stmt.(type) {
	case *sqlparser.Delete:
		tableNameAliasLst = parseSQLTableExprs(stmt.TableExprs)
		if stmt.Where != nil {
			var ok bool
			filterFields, ok = parseSQLWhere(db, stmt.Where, tableNameAliasLst)
			if !ok {
				return nil, nil, false
			}
		}

	default:
		logrus.Fatalf("[APP SQL PARSER] Unsupported SQL statement for delete parser: %s", stmtStr)
	}

	tableNames := make([]string, len(tableNameAliasLst))
	for i, tableAlias := range tableNameAliasLst {
		tableNames[i] = tableAlias.name
	}

	return filterFields, tableNames, true
}

func parseSQLTableExprs(tableExprs sqlparser.TableExprs) []tableNameAlias {
	var tableNameAliasLst []tableNameAlias
	for _, table := range tableExprs {
		if aliasedTableExpr, ok := table.(*sqlparser.AliasedTableExpr); ok {
			if tableName, ok := aliasedTableExpr.Expr.(sqlparser.TableName); ok {
				tableNameAliasLst = append(tableNameAliasLst, tableNameAlias{alias: aliasedTableExpr.As.CompliantName(), name: tableName.Name.CompliantName()})
			}
		}
	}
	return tableNameAliasLst
}

func parseColumnName(compliantName string) (string, string) {
	splits := strings.SplitAfterN(compliantName, ".", 1)
	if len(splits) == 2 {
		return splits[0], splits[1]
	}
	return "", compliantName
}

func parseTableName(prefixTableName string, tableNameAliasLst []tableNameAlias) string {
	if prefixTableName == "" {
		return tableNameAliasLst[0].name
	}
	for _, t := range tableNameAliasLst {
		if t.alias != "" && prefixTableName == t.alias {
			return t.name
		}
		if prefixTableName == t.name {
			return t.name
		}
	}
	logrus.Fatalf("unexpected")
	return tableNameAliasLst[0].name
}
