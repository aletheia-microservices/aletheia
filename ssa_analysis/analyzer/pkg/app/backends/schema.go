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
	database    *Database
}

func NewSchema(name string, database *Database) *Schema {
	return &Schema{
		name:     name,
		fields:   make(map[string]*Field),
		database: database,
	}
}

func (schema *Schema) getConstraintsForeignKey() []*Constraint {
	var constraints []*Constraint
	for _, constraint := range schema.constraints {
		if constraint.IsForeignKey() {
			constraints = append(constraints, constraint)
		}
	}
	return constraints
}

func (schema *Schema) GetAllFieldsLst() []*Field {
	fieldsLst := make([]*Field, len(schema.fields))
	i := 0
	for _, field := range schema.fields {
		fieldsLst[i] = field
		i++
	}
	return fieldsLst
}

func (schema *Schema) GetName() string {
	return schema.name
}

func (schema *Schema) GetDatabase() *Database {
	return schema.database
}

func (schema *Schema) MarshalJSON() ([]byte, error) {
	fieldsLst := make([]string, len(schema.fields))
	i := 0
	for _, field := range schema.fields {
		fieldsLst[i] = field.GetPath()
		i++
	}

	// sort fields
	sort.Slice(fieldsLst, func(i, j int) bool {
		return fieldsLst[i] < fieldsLst[j]
	})

	constraintsLst := make([]string, len(schema.constraints))
	i = 0
	for _, constraint := range schema.constraints {
		constraintsLst[i] = constraint.String()
		i++
	}

	// sort constraints
	sort.Strings(constraintsLst)

	return json.Marshal(&struct {
		Name        string   `json:"name"`
		Fields      []string `json:"fields"`
		Constraints []string `json:"constraints"`
	}{
		Name:        schema.name,
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
		log.Panicf("field with path (%s) not found for schema: %s\n", path, schema.String())
	}
	return field
}

func (schema *Schema) GetFieldByPathIfExists(path string) *Field {
	field, ok := schema.fields[path]
	if ok {
		return field
	}
	return nil
}

func (schema *Schema) GetOrCreateField(database *Database, path string) *Field {
	field, ok := schema.fields[path]
	if !ok {
		field = NewField(path, database, schema)
		schema.AddField(field)
	}
	return field
}

func (schema *Schema) RemoveConstraint(old *Constraint) {
	for i, c := range schema.constraints {
		if c == old {
			schema.constraints = append(schema.constraints[:i], schema.constraints[i+1:]...)
			return
		}
	}
}

func (schema *Schema) AddConstraint(newConstraint *Constraint) {
	for _, existingConstraint := range schema.getConstraintsForeignKey() {
		if existingConstraint.GetFieldAt(0) == newConstraint.GetFieldAt(0) &&
			existingConstraint.GetFieldAt(1) == newConstraint.GetFieldAt(1) {

			if existingConstraint.mandatory == newConstraint.mandatory {
				// ignore if constraint already exists
				return
			} else {
				for reqIdx, mandatory := range newConstraint.reqIdxToMandatory {
					if ok, m := existingConstraint.reqIdxToMandatory[reqIdx]; ok {
						// upgrade existing constraint to mandatory set to true
						if !m {
							existingConstraint.reqIdxToMandatory[reqIdx] = mandatory
						}
					} else {
						// add new entry with mandatory set to true
						existingConstraint.reqIdxToMandatory[reqIdx] = mandatory
					}
				}
				return
			}
		}
	}
	schema.constraints = append(schema.constraints, newConstraint)
}

func (schema *Schema) GetAllConstraints() []*Constraint {
	return schema.constraints
}

func (schema *Schema) GetAllForeignKeyConstraints() []*Constraint {
	var constraints []*Constraint
	for _, constraint := range schema.constraints {
		if constraint.IsForeignKey() {
			constraints = append(constraints, constraint)
		}
	}
	return constraints
}
