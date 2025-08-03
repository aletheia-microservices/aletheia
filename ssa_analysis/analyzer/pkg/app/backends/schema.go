package backends

import (
	"encoding/json"
	"log"
	"sort"
)

type Schema struct {
	name        string // can be name of (sql) table, (nosql) collection, or (queue) topic
	fields      map[string]*Field
	constraints []*Constraint
}

func NewSchema(name string) *Schema {
	return &Schema{
		fields: make(map[string]*Field),
	}
}

func (schema *Schema) GetName() string {
	return schema.name
}

func (schema *Schema) MarshalJSON() ([]byte, error) {
	fieldsLst := make([]string, len(schema.fields))
	i := 0
	for _, field := range schema.fields {
		fieldsLst[i] = field.GetPath()
		i++
	}

	// sort by field path
	sort.Slice(fieldsLst, func(i, j int) bool {
		return fieldsLst[i] < fieldsLst[j]
	})

	constraintsLst := make([]string, len(schema.constraints))
	i = 0
	for _, constraint := range schema.constraints {
		constraintsLst[i] = constraint.String()
		i++
	}

	return json.Marshal(&struct {
		Fields      []string `json:"fields"`
		Constraints []string `json:"constraints"`
	}{
		Fields:      fieldsLst,
		Constraints: constraintsLst,
	})
}

func (schema *Schema) String() string {
	fieldsLst := make([]string, len(schema.fields))
	i := 0
	for _, field := range schema.fields {
		fieldsLst[i] = field.GetPath()
		i++
	}
	sort.Strings(fieldsLst)

	var str string
	for _, field := range fieldsLst {
		str += "\t " + field + "\n"
	}
	for _, constraint := range schema.constraints {
		str += "\t " + constraint.String() + "\n"
	}
	return str
}

func (schema *Schema) HasField(fieldname string) bool {
	_, ok := schema.fields[fieldname]
	return ok
}

func (schema *Schema) AddField(field *Field) {
	schema.fields[field.path] = field
}

func (schema *Schema) GetFields() map[string]*Field {
	return schema.fields
}

func (schema *Schema) GetFieldByPath(path string) *Field {
	field, ok := schema.fields[path]
	if !ok {
		log.Fatalf("field with path (%s) not found for schema: %s\n", path, schema.String())
	}
	return field
}

func (schema *Schema) GetOrCreateField(database *Database, path string) *Field {
	field, ok := schema.fields[path]
	if !ok {
		field = NewField(path, database)
		schema.fields[path] = field
	}
	return field
}

func (schema *Schema) AddConstraint(constraint *Constraint) {
	schema.constraints = append(schema.constraints, constraint)
}
