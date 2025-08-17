package app

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"strings"

	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/auxten/postgresql-parser/pkg/walk"
	"github.com/xwb1989/sqlparser"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/app/services"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/utils"
)

func (app *App) Init() {
	servicesInfo, datastoresInfo, frontends := blueprint.LoadWiring(app.GetName())

	// parse services
	for _, svcInfo := range servicesInfo {
		name := svcInfo.Name
		constructor := svcInfo.ConstructorName
		pkg := svcInfo.Package
		pkgpath := svcInfo.PackagePath
		path := svcInfo.PackagePath + "." + svcInfo.Name
		impl := svcInfo.Name + "Impl"
		methods := svcInfo.Methods
		args := svcInfo.ServiceArgs

		service := services.NewService(name, impl, pkg, pkgpath, path, constructor, args)
		service.SetMethods(methods...)
		app.AddService(service)
	}

	for _, svcInfo := range servicesInfo {
		log.Printf("service = %s, args = %v\n", svcInfo.Name, svcInfo.ServiceArgs)
	}

	for _, svcInfo := range servicesInfo {
		service := app.GetServiceByName(svcInfo.Name)
		for _, dep := range svcInfo.Edges {
			otherService := app.GetServiceByName(dep)
			service.AddDependency(otherService)
		}
	}

	// parse databases
	for _, dsInfo := range datastoresInfo {
		database := backends.NewDatabase(dsInfo.Name, dsInfo.GetTypeString())
		app.AddDatabase(database)
	}

	// parse entrypoints
	for _, serviceName := range frontends {
		service := app.GetServiceByName(serviceName)
		for _, method := range service.GetMethods() {
			app.AddEntrypoint(service, method)
		}
	}
	for _, service := range app.GetAllServices() {
		if service.HasInitializerMethod() {
			// Run() method can also be considered as entrypoint
			// because they are always called when initializing services
			app.AddEntrypoint(service, "Run")
		}
	}
}

func (app *App) InitServiceFields(pkgs []*ssa.Package) {
	for _, pkg := range pkgs {
		fmt.Printf("[APP PARSER] analyzing package: %s\n", pkg.String())
		for _, member := range pkg.Members {
			if ssaType, ok := member.(*ssa.Type); ok {
				service := app.GetServiceWithImplPathIfExists(ssaType.String())
				if service == nil {
					continue
				}
				fmt.Printf("\t[APP PARSER] found service impl: %s\n", service.GetImpl())
				if typeNamed, ok := ssaType.Type().(*types.Named); ok {
					if typeStruct, ok := typeNamed.Underlying().(*types.Struct); ok {
						i := 0
						for i < typeStruct.NumFields() {
							typeVar := typeStruct.Field(i)
							field := services.NewField(i, typeVar.Name())
							service.AddField(field)
							fmt.Printf("\t\t[APP PARSER] created new field: %s\n", field.String())
							i++
						}
					}
				}
			}
		}
	}
}

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

func (app *App) ParseSchemaFromUserFile() {
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
			log.Fatalf("error reading sql files: %s", err.Error())
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
	fmt.Printf("[APP SQL PARSER] parsing statement: %s\n", sql)

	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	sql = strings.TrimSpace(sql)

	var tableName string
	var fields map[string]*backends.Field

	w := &walk.AstWalker{
		Fn: func(ctx interface{}, node interface{}) (stop bool) {
			fmt.Printf("[APP SQL PARSER] visiting node (%T): %v\n", node, node)

			switch stmt := node.(type) {
			case *tree.CreateTable:
				tableName = stmt.Table.Table()
				fields = make(map[string]*backends.Field, 0)

			case *tree.ColumnTableDef:
				columnName := stmt.Name.String()
				fieldName := tableName + "." + columnName
				fieldPath := database.GetName() + "." + fieldName

				schema := database.GetSchemaByNameIfExists(tableName)
				if schema == nil {
					schema = backends.NewSchema(tableName)
					database.AddSchema(schema)
				}

				field := schema.GetFieldByPathIfExists(fieldPath)
				if field == nil {
					field = backends.NewField(fieldPath, database, schema)
					schema.AddField(field)
					fmt.Printf("[APP SQL PARSER] added new database field: %s\n", field.GetPath())
				}
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
					fmt.Printf("[APP SQL PARSER] added new constraint: %s\n", constraint.String())
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
				fmt.Printf("[APP SQL PARSER] added new constraint: %s\n", constraint.String())

			}
			return false
		},
	}

	stmts, err := parser.Parse(sql)
	if err != nil {
		log.Fatalf("[APP SQL PARSER] %s", err.Error())
		return
	}

	ok, err := w.Walk(stmts, nil)
	if err != nil {
		log.Fatalf("[APP SQL PARSER] %s", err.Error())
	} else if !ok {
		log.Fatalf("[APP SQL PARSER] UNEXPECTED!")
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
func ParseSQLRead(db string, stmtStr string) ([]string, []string, []string) {
	fmt.Printf("[APP SQL PARSER] parsing sql read for database (%s): %s\n", db, stmtStr)
	stmt, err := sqlparser.Parse(stmtStr)
	if err != nil {
		log.Fatalf("[APP SQL PARSER] unable to parse sql query (%s): %s", stmtStr, err.Error())
	}

	//argIdx := 0
	//FIXME: selectedFields has always only one element and the object is the "target" from the func args
	var selectedFields []string
	var filterFields []string
	var tableNameAliasLst []tableNameAlias

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		tableNameAliasLst = parseSQLTableExprs(stmt.From)

		var readAllFields bool
		for _, expr := range stmt.SelectExprs {
			if _, ok := expr.(*sqlparser.StarExpr); ok {
				readAllFields = true
				fieldpath := db + "." + tableNameAliasLst[0].name // + ".*"
				selectedFields = append(selectedFields, fieldpath)
				fmt.Printf("[APP SQL PARSER] found sqlparser.StarExpr (%t)\n", readAllFields)
			} else if aliasedExpr, ok := expr.(*sqlparser.AliasedExpr); ok {
				if valTuple, ok := aliasedExpr.Expr.(sqlparser.ValTuple); ok {
					for rowIdx, expr := range valTuple {
						if col, ok := expr.(*sqlparser.ColName); ok {
							prefixTableName, columnName := parseColumnName(string(col.Name.CompliantName()))
							tableName := parseTableName(prefixTableName, tableNameAliasLst)
							fieldpath := db + "." + tableName + "." + columnName

							selectedFields = append(selectedFields, fieldpath)
							fmt.Printf("[APP SQL PARSER] [SELECT record %d/%d]: %s\n", rowIdx+1, len(valTuple), fieldpath)
						}
					}
				}
			}
		}

		filterFields = parseSQLWhere(db, stmt.Where, tableNameAliasLst)

	default:
		log.Fatalf("[APP SQL PARSER] Unsupported SQL statement: %s", stmtStr)
	}

	tableNames := make([]string, len(tableNameAliasLst))
	for i, tableAlias := range tableNameAliasLst {
		tableNames[i] = tableAlias.name
	}

	return selectedFields, filterFields, tableNames
}

func parseSQLWhere(db string, stmtWhere *sqlparser.Where, tableNameAliasLst []tableNameAlias) []string {
	var filterFields []string
	if comparisonExpr, ok := stmtWhere.Expr.(*sqlparser.ComparisonExpr); ok {
		var leftFieldName string
		if col, ok := comparisonExpr.Left.(*sqlparser.ColName); ok {
			prefixTableName, columnName := parseColumnName(string(col.Name.CompliantName()))
			tableName := parseTableName(prefixTableName, tableNameAliasLst)
			leftFieldName = db + "." + tableName + "." + columnName
		}
		filterFields = append(filterFields, leftFieldName)
		fmt.Printf("[APP SQL PARSER] [WHERE]: %s -> <SOME OBJECT TBD>\n", leftFieldName)
	}
	return filterFields
}

// ParseSQLWrite returns two slices, one for written fields and another for filter fields
// where each field return has format <table>.<field>
func ParseSQLWrite(db string, stmtStr string) ([]string, []string, string) {
	fmt.Printf("[APP SQL PARSER] parsing sql write for database (%s): %s\n", db, stmtStr)

	stmt, err := sqlparser.Parse(stmtStr)
	if err != nil {
		log.Fatalf("[APP SQL PARSER] unable to parse sql query (%s): %s", stmtStr, err.Error())
	}

	var writtenFields []string
	var filterFields []string
	var tableName string

	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		if values, ok := stmt.Rows.(sqlparser.Values); ok {
			for rowIdx, tuple := range values {
				for colIdx, expr := range tuple {
					if sqlVal, ok := expr.(*sqlparser.SQLVal); ok {
						tableName = stmt.Table.Name.CompliantName()
						if sqlVal.Type == sqlparser.ValArg { // placeholder (e.g., '?' that is then parsed into ':v1', ':v2', etc.)
							fieldpath := db + "." + tableName + "." + stmt.Columns[colIdx].CompliantName()
							placeholderVal := string(sqlVal.Val)
							writtenFields = append(writtenFields, fieldpath)
							fmt.Printf("[APP SQL PARSER] (record %d/%d) INSERT %s = (%s) -> <SOME OBJECT TBD>\n", rowIdx+1, len(values), fieldpath, placeholderVal)
						}
					}
				}
			}
		} else {
			log.Fatalf("[APP SQL PARSER] unexpected type %T for rows in sql insert: %s", stmt.Rows, stmtStr)
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
					placeholderVal := string(sqlVal.Val)
					writtenFields = append(writtenFields, fieldpath)
					fmt.Printf("[APP SQL PARSER] SET %s = (%s) -> <SOME OBJECT TBD>\n", fieldpath, placeholderVal)
					argIdx++
				}
			}
		}

		filterFields = parseSQLWhere(db, stmt.Where, tableNameAliasLst)
	}
	return writtenFields, filterFields, tableName
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
	log.Fatal("unexpected")
	return tableNameAliasLst[0].name
}
