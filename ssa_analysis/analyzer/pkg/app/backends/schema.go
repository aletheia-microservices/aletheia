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
	// index: (field0, field1) -> constraint
	fkIndex map[[2]*Field]*Constraint
	// no need to remove if referenced is deleted!
	refBy   map[*Schema]bool
}

func NewSchema(name string, database *Database) *Schema {
	return &Schema{
		name:     name,
		fields:   make(map[string]*Field),
		database: database,
		fkIndex:  make(map[[2]*Field]*Constraint),
		refBy:    make(map[*Schema]bool),
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
			if old.IsForeignKey() {
				delete(schema.fkIndex, [2]*Field{old.GetFieldAt(0), old.GetFieldAt(1)})
			}
			return
		}
	}
}

func (schema *Schema) GetForeignKeyForPair(f1 *Field, f2 *Field) *Constraint {
	key := [2]*Field{f1, f2}
	if existing, ok := schema.fkIndex[key]; ok {
		return existing
	}
	return nil
}

func (schema *Schema) AddConstraint(newConstraint *Constraint) bool {
	if newConstraint.IsForeignKey() {
		key := [2]*Field{
			newConstraint.GetFieldAt(0),
			newConstraint.GetFieldAt(1),
		}

		if existing, ok := schema.fkIndex[key]; ok {
			// same fields, merge semantics
			if existing.mandatory == newConstraint.mandatory {
				// identical, nothing to do
				return false
			}

			for reqIdx, mandatory := range newConstraint.reqIdxToMandatory {
				if cur, ok := existing.reqIdxToMandatory[reqIdx]; !ok || (!cur && mandatory) {
					// either new index or upgrade false -> true
					existing.reqIdxToMandatory[reqIdx] = mandatory
				}
			}
			return false
		}
		schema.fkIndex[key] = newConstraint
		newConstraint.GetFieldAt(1).GetSchema().refBy[schema] = true
	}

	schema.constraints = append(schema.constraints, newConstraint)
	return true
}

func (schema *Schema) AddSchemaToRefBy(other *Schema) {
	schema.refBy[other] = true
}

func (schema *Schema) GetAllSchemasRefBy() map[*Schema]bool {
	return schema.refBy
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
