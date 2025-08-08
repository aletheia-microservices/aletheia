package backends

import "strings"

type Field struct {
	path        string
	database    *Database
	schema      *Schema
	constraints []*Constraint
}

func NewField(path string, database *Database, schema *Schema) *Field {
	return &Field{
		path:     path,
		database: database,
		schema:   schema,
	}
}

func (field *Field) GetPath() string {
	return field.path
}

// extract <name> from <db>.<schema>.<name
func (field *Field) GetName() string {
	if idx := strings.LastIndex(field.path, "."); idx != -1 {
		return field.path[idx+1:]
	}
	return field.path
}

func (field *Field) GetSchema() *Schema {
	return field.schema
}

func (field *Field) GetDatabase() *Database {
	return field.database
}

func (field *Field) AddConstraint(constraint *Constraint) {
	field.constraints = append(field.constraints, constraint)
}

func (field *Field) HasConstraintForeignKeyToField(otherField *Field) bool {
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_FOREIGN_KEY && constraint.fields[1] == otherField {
			return true
		}
	}
	return false
}

// searches for unicity in single field
// also includes primary key
func (field *Field) HasContraintUnicity() bool {
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_UNIQUE || constraint.t == CONSTRAINT_PRIMARY {
			if len(constraint.fields) == 1 {
				return true
			}
		}
	}
	return false
}

func (field *Field) String() string {
	str := field.path
	if len(field.constraints) > 0 {
		str += "\n"
	}
	for i, constraint := range field.constraints {
		str += "\t\t" + constraint.String()
		if i < len(constraint.fields) {
			str += "\n"
		}
	}
	return str
}
