package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"analyzer/pkg/app/backends"
)

type ForeignKey struct {
	srcDatabase, srcTable, srcField string
	dstDatabase, dstTable, dstField string
}

func (app *App) ParseUserInputReferences() {
	dir := filepath.Join("input", app.GetName())

	entries, err := os.ReadDir(dir)
	if err != nil {
		logrus.Panicf("[ERROR] failed to read dir (%s): %s", dir, err.Error())
	}

	var foreignkeys []ForeignKey

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".yml" && ext != ".yaml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			logrus.Panicf("[ERROR] failed to read (%s): %s", path, err.Error())
		}

		var lines []string
		if err := yaml.Unmarshal(data, &lines); err != nil {
			logrus.Panicf("[ERROR] failed to unmarshal (%s): %s", path, err.Error())
		}

		for _, line := range lines {
			foreignkey := parseReferenceEntry(line)
			if err != nil {
				continue
			}
			foreignkeys = append(foreignkeys, *foreignkey)
		}
	}

	for _, foreignkey := range foreignkeys {
		database1 := app.GetDatabaseByName(foreignkey.srcDatabase)
		schema1 := database1.GetOrCreateSchema(foreignkey.srcTable)
		fieldpath1 := foreignkey.srcDatabase + "." + foreignkey.srcTable + "." + foreignkey.srcField
		field1 := schema1.GetOrCreateField(database1, fieldpath1)

		database2 := app.GetDatabaseByName(foreignkey.dstDatabase)
		schema2 := database2.GetOrCreateSchema(foreignkey.dstTable)
		fieldpath2 := foreignkey.dstDatabase + "." + foreignkey.dstTable + "." + foreignkey.dstField
		field2 := schema2.GetOrCreateField(database2, fieldpath2)

		constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field1, field2)
		field1.AddConstraint(constraint)
		schema1.AddConstraint(constraint)
	}
}

func parseReferenceEntry(line string) *ForeignKey {
	// e.g., FOREIGN_KEY delivery_db.delivery.StationName REFERENCES station_db.station.Name

	parts := strings.Fields(line)
	if len(parts) != 4 || parts[0] != "FOREIGN_KEY" || parts[2] != "REFERENCES" {
		logrus.Panicf("[ERROR] unexpected foreign key format: %s", line)
	}

	src := parts[1]
	dst := parts[3]

	parseTriple := func(s string) (db, table, col string) {
		p := strings.Split(s, ".")
		if len(p) != 3 {
			logrus.Panicf("[ERROR] expected db.table.column, got: %s", s)
		}
		return p[0], p[1], p[2]
	}

	srcDB, srcTable, srcField := parseTriple(src)
	dstDB, dstTable, dstField := parseTriple(dst)

	return &ForeignKey{
		srcDatabase: srcDB,
		srcTable:    srcTable,
		srcField:    srcField,
		dstDatabase: dstDB,
		dstTable:    dstTable,
		dstField:    dstField,
	}
}
