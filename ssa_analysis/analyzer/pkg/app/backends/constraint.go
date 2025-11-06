package backends

import "log"

type ConstraintType int

const (
	CONSTRAINT_FOREIGN_KEY ConstraintType = iota
	CONSTRAINT_UNIQUE      ConstraintType = iota
	CONSTRAINT_PRIMARY     ConstraintType = iota
)

type Constraint struct {
	t ConstraintType
	// for foreign key constraint, index 0 is for field referencing and index 1 is for field being referenced
	fields    []*Field
	mandatory bool // for foreign key constraints

	// 1. applies to foreign key only
	// 2. the entire constraint is non-mandatory if any keyvalue is non-mandatory
	// 3. for a given request index:
	// - when a given execution is mandatory, the mandatory flag is always overwritten to true
	// - when a given execution is non-mandatory, the mandatory flag is only set to false if it is not yet true
	// key is end-to-end request index
	// map is mandatory bool
	reqIdxToMandatory map[int]bool
}

func NewConstraint(t ConstraintType, fields ...*Field) *Constraint {
	return &Constraint{
		t:                 t,
		fields:            fields,
		reqIdxToMandatory: make(map[int]bool),
	}
}

func (constraint *Constraint) IsMandatory() bool {
	for _, m := range constraint.reqIdxToMandatory {
		if !m {
			return false
		}
	}
	return true
}

func (constraint *Constraint) GetRequestsIndexesOnMandatoryFlags() []int {
	var reqIdxs []int
	for idx, m := range constraint.reqIdxToMandatory {
		if m {
			reqIdxs = append(reqIdxs, idx)
		}
	}
	return reqIdxs
}

func (constraint *Constraint) EnableMandatory(reqIdx int) bool {
	constraint.reqIdxToMandatory[reqIdx] = true
	constraint.mandatory = true
	return true
}

func (constraint *Constraint) DisableMandatory(reqIdx int) bool {
	if m, ok := constraint.reqIdxToMandatory[reqIdx]; ok {
		if m {
			return false
		}
		constraint.reqIdxToMandatory[reqIdx] = false
		return true
	}

	constraint.reqIdxToMandatory[reqIdx] = false
	constraint.mandatory = false
	return true
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
	if idx > len(constraint.fields) - 1 {
		log.Panicf("[CONSTRAINT] index (%d) out of bounds for constraint: %v\n", idx, constraint.String())
	}
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
