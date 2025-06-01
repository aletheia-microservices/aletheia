package unicity

import (
	"slices"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
)

type Operation struct {
	idx                 int
	call                *abstractgraph.AbstractDatabaseCall
	datastore           *datastores.Datastore
	constraints         []*datastores.Constraint
	writtenFields       []*datastores.Field
	onUnicityConstraint bool
	/* affectedOps         []*Operation */
	// key: following ops that are affected by current constraint
	// value: fields in the current op that were referenced in the following op
	affectedOps map[*Operation][]*datastores.Field
}

func (op *Operation) HasWrittenField(f *datastores.Field) bool {
	return slices.Contains(op.writtenFields, f)
}

/* func (op *Operation) AddAffectedOp(affectedOp *Operation) {
	op.affectedOps = append(op.affectedOps, affectedOp)
} */

func (op *Operation) AddAffectedOpAndReferencedField(newAffectedOp *Operation, refField *datastores.Field) {
	if op.affectedOps == nil {
		op.affectedOps = make(map[*Operation][]*datastores.Field)
	}
	op.affectedOps[newAffectedOp] = append(op.affectedOps[newAffectedOp], refField)
}

func (op *Operation) AffectsOps() bool {
	return len(op.affectedOps) > 0
}

func (op *Operation) GetAffectedOps() map[*Operation][]*datastores.Field {
	return op.affectedOps
}

func NewOperation(idx int, call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore, writtenFields []*datastores.Field, constraints []*datastores.Constraint) *Operation {
	return &Operation{
		idx:                 idx,
		call:                call,
		datastore:           datastore,
		writtenFields:       writtenFields,
		constraints:         constraints,
		onUnicityConstraint: len(constraints) > 0,
	}
}
