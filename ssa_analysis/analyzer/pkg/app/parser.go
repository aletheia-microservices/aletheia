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
func parseSQLStatementsFromInput(input string) []SQLDbStmt {
	var dbStmts []SQLDbStmt
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
			dbStmts = append(dbStmts, SQLDbStmt{db, stmt})
		}
	}
	return dbStmts
}

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

type SQLDbStmt struct {
	db   string
	stmt string
}

func parseSQLStatement(database *backends.Database, sql string) {
	fmt.Printf("[SQL PARSER] parsing statement: %s\n", sql)

	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	sql = strings.TrimSpace(sql)

	var tableName string
	var fields map[string]*backends.Field

	w := &walk.AstWalker{
		Fn: func(ctx interface{}, node interface{}) (stop bool) {
			fmt.Printf("[SQL PARSER] visiting node (%T): %v\n", node, node)

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
					fmt.Printf("[SQL PARSER] added new database field: %s\n", field.GetPath())
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
					fmt.Printf("[SQL PARSER] added new constraint: %s\n", constraint.String())
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
				fmt.Printf("[SQL PARSER] added new constraint: %s\n", constraint.String())

			}
			return false
		},
	}

	stmts, err := parser.Parse(sql)
	if err != nil {
		log.Fatalf("[SQL PARSER] %s", err.Error())
		return
	}

	ok, err := w.Walk(stmts, nil)
	if err != nil {
		log.Fatalf("[SQL PARSER] %s", err.Error())
	} else if !ok {
		log.Fatalf("[SQL PARSER] UNEXPECTED!")
	}
}
