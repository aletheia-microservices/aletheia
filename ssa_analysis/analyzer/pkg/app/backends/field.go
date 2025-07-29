package backends

type Field struct {
	path        string
	database    *Database
	constraints []*Constraint
}

func NewField(path string, database *Database) *Field {
	return &Field{
		path:     path,
		database: database,
	}
}

func (field *Field) GetPath() string {
	return field.path
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
