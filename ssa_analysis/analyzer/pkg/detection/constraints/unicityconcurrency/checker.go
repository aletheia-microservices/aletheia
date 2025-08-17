package unicityconcurrency

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/utils"
)

type VulnerableWriteSet struct {
	constrainedOp     *WriteOperation
	otherOps          []*WriteOperation
	constrainedFields []*backends.Field
}

func (writeSet *VulnerableWriteSet) addOtherOperation(op *WriteOperation) {
	writeSet.otherOps = append(writeSet.otherOps, op)
}

func (writeSet *VulnerableWriteSet) hasOtherOperation(op *WriteOperation) bool {
	return slices.Contains(writeSet.otherOps, op)
}

func (detector *UnicityConcurrencyDetector) checkInconsistency(app *app.App, request *Request, currOp *WriteOperation) {
	dbname := currOp.call.GetToNode().GetDatabaseName()
	db := app.GetDatabaseByName(dbname)

	var constrainedFields []*backends.Field
	for _, arg := range currOp.arguments {
		fmt.Printf("[UNICITY CHECKER] arg (%s) has primary taint lst: %v\n", arg.String(), arg.GetPrimaryTaintsFlatList())
		for _, taint := range arg.GetPrimaryTaintsFlatList() {
			fieldpath := taint.GetDatabasePath()

			// [TO BE IMPROVED]
			// there may be cases where primary taint is not related to this database
			// when services make more than one call to different databases
			//
			// in the future, we may just associate the taint with the call ID
			// and then just check if the IDs match
			if dbname == utils.ExtractDatabaseNameFromFieldPath(fieldpath) {
				schema := db.GetSchemaByNameIfExists(utils.ExtractSchemaNameFromFieldPath(fieldpath))
				if schema == nil {
					schemaName := utils.ExtractSchemaNameFromFieldPath(fieldpath)
					if strings.HasSuffix(schemaName, "[*]") {
						// [TO BE IMPROVED]
						// sometimes we get schema name "schema[*]" from fieldpaths "schema[*].Value"
						// because mongodb read filter fields are not being yet parsed for reads taints
						// for now we hardcode to remove the [*] in "schema[*]"
						schemaName = schemaName[:len(schemaName)-3]
						schema = db.GetSchemaByNameIfExists(schemaName)
					} else {
						log.Panicf("[UNICITY CHECKER] nil schema (%s) for fieldpath (%s)\n", utils.ExtractSchemaNameFromFieldPath(fieldpath), fieldpath)
					}
				}
				fmt.Printf("[UNICITY CHECKER] get field for (%s) in schema (%s)\n", fieldpath, schema.GetName())
				field := schema.GetFieldByPath(fieldpath)

				// [TO BE IMPROVED]
				// in the future, the ssa parser should be the one to infer 
				// all schema fields (from AST structure) beforehand
				if field == nil {
					field = backends.NewField(fieldpath, db, schema)
					schema.AddField(field)
				}

				if field.HasContraintUnicity() && !slices.Contains(constrainedFields, field) {
					constrainedFields = append(constrainedFields, field)
				}
			}
		}
	}

	// same logic as in foreignkeycoordination and foreignkeycascade
	// but here we verify if secondaryTaint.IsWrite()
	for _, arg := range currOp.arguments {
		for _, secondaryTaint := range arg.GetSecondaryTaintsFlatList() {
			if secondaryTaint.GetDatabaseCallID() != currOp.GetCallID() && secondaryTaint.IsWrite() {
				otherOp := request.FindOperationByCallID(secondaryTaint.GetDatabaseCallID())
				if otherOp != nil {
					otherWriteSet := detector.findVulnerableWriteSetForOperation(request, otherOp)
					if otherWriteSet != nil && !otherWriteSet.hasOtherOperation(currOp) {
						otherWriteSet.addOtherOperation(currOp)
					}
				}
			}
		}
	}

	if constrainedFields != nil {
		writeSet := &VulnerableWriteSet{
			constrainedOp:     currOp,
			constrainedFields: constrainedFields,
		}
		detector.addVulnerableWriteSet(request, writeSet)
	}
}
