package foreign_key

import (
	"fmt"

	"analyzer/pkg/datastores"
	"analyzer/pkg/types"
)

type ForeignKeyRead struct {
	refField     *datastores.Entry // field that is referencing
	originField  *datastores.Entry // field that is being referenced
	refDbCall    *types.ParsedDatabaseCall
	originDbCall *types.ParsedDatabaseCall
}

func newForeignKeyRead(refField *datastores.Entry, originField *datastores.Entry, refDbCall *types.ParsedDatabaseCall, originDbCall *types.ParsedDatabaseCall) *ForeignKeyRead {
	return &ForeignKeyRead{
		refField:     refField,
		originField:  originField,
		refDbCall:    refDbCall,
		originDbCall: originDbCall,
	}
}

func (read *ForeignKeyRead) String() string {
	ref := fmt.Sprintf("- ref:\t%s\n\t@ %s", read.refField.GetFullName(), read.refDbCall.String())
	dst := fmt.Sprintf("- dst:\t%s\n\t@ %s", read.originField.GetFullName(), read.originDbCall.String())
	return ref + "\n" + dst
}
