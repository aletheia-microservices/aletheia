package backends

type ConstraintType int

const (
	CONSTRAINT_FOREIGN_KEY ConstraintType = iota
	CONSTRAINT_UNIQUE      ConstraintType = iota
	CONSTRAINT_PRIMARY     ConstraintType = iota
)

type Constraint struct {
	t ConstraintType
	// for foreign key constraint, index 0 is for field referencing and index 1 is for field being referenced
	fields []*Field
	mandatory bool // for foreign key constraints
}

func NewConstraint(t ConstraintType, fields ...*Field) *Constraint {
	return &Constraint{
		t:      t,
		fields: fields,
	}
}

func (constraint *Constraint) IsMandatory() bool {
	return constraint.mandatory
}

func (constraint *Constraint) EnableMandatory() {
	constraint.mandatory = true
}

func (constraint *Constraint) DisableMandatory() {
	constraint.mandatory = false
}

func (constraint *Constraint) IsForeignKey() bool {
	return constraint.t == CONSTRAINT_FOREIGN_KEY
}

// also includes primary keys
func (constraint *Constraint) IsUnique() bool {
	return constraint.t == CONSTRAINT_UNIQUE || constraint.t == CONSTRAINT_PRIMARY
}

func (constraint *Constraint) AddField(field *Field) {
	constraint.fields = append(constraint.fields, field)
}

func (constraint *Constraint) GetFields() []*Field {
	return constraint.fields
}

func (constraint *Constraint) GetFieldAt(idx int) *Field {
	return constraint.fields[idx]
}

func (constraint *Constraint) GetFieldsNamesString() string {
	str := "("
	for i, field := range constraint.fields {
		str += field.GetPath()
		if i < len(constraint.fields)-1 {
			str += ", "
		}
	}
	str += ")"
	return str
}

func (constraint *Constraint) String() string {
	switch constraint.t {
	case CONSTRAINT_FOREIGN_KEY:
		var suffix string
		if constraint.IsMandatory() {
			suffix = " [MANDATORY]"
		}
		return "FOREIGN_KEY " + constraint.fields[0].GetPath() + " REFERENCES " + constraint.fields[1].GetPath() + suffix
	case CONSTRAINT_UNIQUE:
		return "UNIQUE " + constraint.GetFieldsNamesString()
	case CONSTRAINT_PRIMARY:
		return "PRIMARY KEY " + constraint.GetFieldsNamesString()
	}
	return ""
}
