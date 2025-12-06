package specialization

import (
	"fmt"

	"analyzer/pkg/types"
)

type RemovedMandatoryEntity struct {
	dbCall          *types.ParsedDatabaseCall
	mandatoryFields []*mandatoryField
}

func newRemovedMandatoryEntity(dbCall *types.ParsedDatabaseCall, mandatoryFields []*mandatoryField) *RemovedMandatoryEntity {
	return &RemovedMandatoryEntity{
		dbCall:          dbCall,
		mandatoryFields: mandatoryFields,
	}
}

func (rme *RemovedMandatoryEntity) String() string {
	return fmt.Sprintf("%s \n\t@ %s", rme.dbCall.DbInstance.GetDatastore().Schema.GetRootUnfoldedField().GetFullName(), rme.dbCall.String())
}
