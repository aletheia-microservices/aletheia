package backends

import (
	"log"
)

type ConstraintType int

const (
	CONSTRAINT_FOREIGN_KEY ConstraintType = iota
	CONSTRAINT_UNIQUE      ConstraintType = iota
)

type Constraint struct {
	t ConstraintType
	// for foreign key constraint, index 0 is for field referencing and index 1 is for field being referenced
	fields []*Field
}

func NewConstraint(t ConstraintType, fields ...*Field) *Constraint {
	return &Constraint{
		t:      t,
		fields: fields,
	}
}

func (constraint *Constraint) GetFields() []*Field {
	return constraint.fields
}

func (constraint *Constraint) String() string {
	if constraint.t == CONSTRAINT_FOREIGN_KEY {
		if constraint.fields[1] == nil {
			log.Fatal("unexpected nil field in index 1 for constraint")
		}
		return "FOREIGN_KEY " + constraint.fields[0].path + " REFERENCES " + constraint.fields[1].path
	}
	return ""
}
