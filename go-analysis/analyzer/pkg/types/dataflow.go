package types

import (
	"analyzer/pkg/types/objects"
)

type Dataflow struct {
	objDataflow *objects.ObjectDataflow
	dbCall      *ParsedDatabaseCall
}

func NewDataflow(objDataflow *objects.ObjectDataflow, dbCall *ParsedDatabaseCall) *Dataflow {
	return &Dataflow{
		objDataflow: objDataflow,
		dbCall:      dbCall,
	}
}

func (df *Dataflow) GetDatabaseCall() *ParsedDatabaseCall {
	return df.dbCall
}

func (df *Dataflow) GetObjectDataflow() *objects.ObjectDataflow {
	return df.objDataflow
}

func (df *Dataflow) HasObjectDataflow(other *objects.ObjectDataflow) bool {
	return df.objDataflow == other
}
