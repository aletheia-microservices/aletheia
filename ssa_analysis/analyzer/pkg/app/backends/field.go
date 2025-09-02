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

func (field *Field) GetConstraints() []*Constraint {
	return field.constraints
}

// extract <name> from <db>.<schema>.<name
func (field *Field) GetName() string {
	parts := strings.SplitN(field.path, ".", 3)
	if len(parts) == 3 {
		return parts[2]
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

func (field *Field) HasConstraintForeignKeyNonMandatoryToField(other *Field) bool {
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_FOREIGN_KEY && !constraint.IsMandatory() && constraint.fields[1] == other {
			return true
		}
	}
	return false
}

func (field *Field) HasConstraintForeignKeyNonMandatory() bool {
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_FOREIGN_KEY && !constraint.IsMandatory() {
			return true
		}
	}
	return false
}

// checks if it is single primary key (not composed)
func (field *Field) IsPrimaryKey() bool {
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_PRIMARY && len(constraint.fields) == 1 {
			return true
		}
	}
	return false
}

// includes primary key
func (field *Field) IsUnique() bool {
	for _, constraint := range field.constraints {
		if (constraint.t == CONSTRAINT_UNIQUE || constraint.t == CONSTRAINT_PRIMARY) && len(constraint.fields) == 1 {
			return true
		}
	}
	return false
}

func (field *Field) GetConstraintForeignKey() []*Constraint {
	var constraints []*Constraint
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_FOREIGN_KEY {
			constraints = append(constraints, constraint)
		}
	}
	return constraints
}

func (field *Field) GetConstraintPrimaryKey() *Constraint {
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_PRIMARY && len(constraint.fields) == 1 {
			return constraint
		}
	}
	return nil
}

func (field *Field) GetConstraintForeignKeyToField(otherField *Field) *Constraint {
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_FOREIGN_KEY && constraint.fields[1] == otherField {
			return constraint
		}
	}
	return nil
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
// EXCLUDES primary key
func (field *Field) HasContraintUnicity() bool {
	for _, constraint := range field.constraints {
		if constraint.t == CONSTRAINT_UNIQUE {
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
