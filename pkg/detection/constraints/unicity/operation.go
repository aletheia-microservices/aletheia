package unicity

import (
	"slices"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
	"analyzer/pkg/types/objects"
)

type Operation struct {
	idx           int
	call          *abstractgraph.AbstractDatabaseCall
	datastore     *datastores.Datastore
	constraints   []*datastores.Constraint
	writtenFields []*datastores.Field
	writtenObjs   []objects.Object
	constrained   bool // constrained by UNIQUE or PK
	affectedOps   map[*Operation][]*datastores.Field
}

func (op *Operation) hasWrittenField(f *datastores.Field) bool {
	return slices.Contains(op.writtenFields, f)
}

func (op *Operation) addAffectedOpAndReferencedField(newAffectedOp *Operation, refField *datastores.Field) {
	if op.affectedOps == nil {
		op.affectedOps = make(map[*Operation][]*datastores.Field)
	}
	op.affectedOps[newAffectedOp] = append(op.affectedOps[newAffectedOp], refField)
}

func (op *Operation) affectsOperations() bool {
	return len(op.affectedOps) > 0
}

func (op *Operation) isConstrained() bool {
	return op.constrained
}

func (op *Operation) getAffectedOperations() map[*Operation][]*datastores.Field {
	return op.affectedOps
}

func (op *Operation) addWrittenObjects(obj ...objects.Object) {
	op.writtenObjs = append(op.writtenObjs, obj...)
}

func (op *Operation) getWrittenObjects() []objects.Object {
	return op.writtenObjs
}

func (op *Operation) getDatastore() *datastores.Datastore {
	return op.datastore
}

func NewOperation(idx int, call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore, writtenFields []*datastores.Field, constraints []*datastores.Constraint) *Operation {
	return &Operation{
		idx:           idx,
		call:          call,
		datastore:     datastore,
		writtenFields: writtenFields,
		constraints:   constraints,
		constrained:   len(constraints) > 0,
	}
}
