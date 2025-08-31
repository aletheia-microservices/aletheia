package keycoordination

import (
	"fmt"
	"slices"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/utils"
)

type ForeignRead struct {
	op1         *ReadOperation // read1 has the secondary taint referencing read2
	op2         *ReadOperation
	field1      *backends.Field
	field2      *backends.Field // field used as foreign key in op2
	constraint1 *backends.Constraint
	constraint2 *backends.Constraint
}

func (detector *KeyCoordinationDetector) checkInconsistency(app *app.App, request *Request) {
	allOps := request.GetAllOperations()
	// check in reverse
	i := len(allOps) - 1
	for i >= 0 {
		detector.checkInconsistencyForOp(app, request, allOps[i])
		i--
	}

}

func (detector *KeyCoordinationDetector) checkInconsistencyForOp(app *app.App, request *Request, currOp *ReadOperation) {
	// same logic as in foreignkeycascade but here we verify if secondaryTaint.IsDelete()
	fmt.Printf("[%s| CHECKER] on op: %v\n", detector.GetTypeStringUpper(), currOp.call.String())
	for _, arg := range currOp.arguments {
		fmt.Printf("[%s | CHECKER] arg (%s)\n", detector.GetTypeStringUpper(), arg.String())
		fmt.Printf("\t[%s | CHECKER] arg (%s) w/ all taints: %v\n", arg.String(), detector.GetTypeStringUpper(), arg.GetAllTaints())

		for objpath, taintLst := range arg.GetAllTaints() {
			var currTaint *abstractgraph.AbstractTaint
			var otherTaints []*abstractgraph.AbstractTaint

			for _, taint := range taintLst {
				if taint.IsPrimary() && taint.GetDatabaseCallID() == currOp.GetCallID() {
					fmt.Printf("\t[%s | CHECKER] objpath={%s} // currTaint={%s}\n", detector.GetTypeStringUpper(), objpath, taint.String())
					currTaint = taint
				} else if !taint.IsPrimary() && taint.GetDatabaseCallID() != currOp.GetCallID() && !slices.Contains(otherTaints, taint) {
					fmt.Printf("\t[%s | CHECKER] objpath={%s} // otherTaint={%s}\n", detector.GetTypeStringUpper(), objpath, taint.String())
					otherTaints = append(otherTaints, taint)
				}
			}

			if currTaint == nil {
				fmt.Printf("\t[%s | CHECKER] skipping currTaint with otherTaints: %v\n", detector.GetTypeStringUpper(), otherTaints)
				continue
			}

			currFieldPath := currTaint.GetDatabasePath()
			currDatabase := app.GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(currFieldPath))
			currField := app.ComputeDatabaseFieldFromPath(currDatabase, currFieldPath)

			for _, otherTaint := range otherTaints {
				fmt.Printf("\t[%s | CHECKER] arg (%s) w/ secondary taint: %v\n", detector.GetTypeStringUpper(), arg.String(), otherTaint.LongString())
				otherOp := request.FindOperationByCallID(otherTaint.GetDatabaseCallID())

				// sanity checks
				if currOp == otherOp || otherOp == nil || currOp.call.GetToNode().GetDatabaseName() == otherOp.call.GetToNode().GetDatabaseName() {
					fmt.Printf("\t[%s | CHECKER] skipping nil op for otherTaint (arg=%s): %s\n", detector.GetTypeStringUpper(), arg.String(), otherTaint.LongString())
					continue
				}

				if !detector.hasForeignRead(request, currOp, otherOp) {
					otherFieldpath := otherTaint.GetDatabasePath()
					otherDatabase := app.GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(otherFieldpath))
					otherField := app.ComputeDatabaseFieldFromPath(otherDatabase, otherFieldpath)

					fmt.Printf("\t\t[%s | CHECKER] currField={%s} // otherField={%s}\n", detector.GetTypeStringUpper(), currField.String(), otherField.String())

					foreignRead := &ForeignRead{
						op1:    currOp,
						op2:    otherOp,
						field1: currField,
						field2: otherField,
					}
					detector.addForeignRead(request, foreignRead)
				}
			}
		}
	}
}

func (detector *KeyCoordinationDetector) updateForeignReadConstraints(fread *ForeignRead) {
	if detector.isTypeForeignKey() {
		fread.constraint1 = fread.field1.GetConstraintForeignKeyToField(fread.field2)
		fread.constraint2 = fread.field2.GetConstraintForeignKeyToField(fread.field1)
	} else if detector.isTypePrimaryKey() {
		fread.constraint1 = fread.field1.GetConstraintPrimaryKey()
		fread.constraint2 = fread.field2.GetConstraintPrimaryKey()
	}
}

// this information is only accurate after the entire schema is built at the end of the iteration
func (detector *KeyCoordinationDetector) isValidForeignRead(fread *ForeignRead) bool {
	if detector.isTypePrimaryKey() && (!fread.field1.IsPrimaryKey() || !fread.field2.IsPrimaryKey()) {
		return false
	}
	if detector.isTypeForeignKey() && fread.field1.IsPrimaryKey() && fread.field2.IsPrimaryKey() {
		return false
	}

	if detector.isTypeForeignKey() && detector.isRestrictive() {
		// 1. restrict detection warnings to mandatory constraints
		// 2. filter out constraints that were created in the current request
		// meaning that this foreign read was either before the write (e.g., to check if exists cache)
		// or after the write, so they cannot be considered as inconsistencies
		//
		// REMINDER: this prevents false positives in e.g., dsb_sn2 on Wrk2APIService.ComposePost():
		// - FOREIGN KEY READ #1:
		// 		READ 1: SocialGraphService.GetFollowers() --> socialgraph_cache.followers.Get()				@ socialgraph_cache.followers.Key
		// 		READ 2: UserTimelineService.WriteUserTimeline() --> usertimeline_db.usertimeline.FindMany() @ usertimeline_db.usertimeline.UserID
		// - FOREIGN KEY READ #2:
		// 		READ 1: SocialGraphService.GetFollowers() --> socialgraph_cache.followers.Get() @ socialgraph_cache.followers.Key
		// 		READ 2: UserTimelineService.WriteUserTimeline() --> usertimeline_cache.*.Get() 	@ usertimeline_cache.*.Key
		// - FOREIGN KEY READ #3:
		// 		READ 1: HomeTimelineService.WriteHomeTimeline() --> hometimeline_cache.*.Get() @ hometimeline_cache.*.Value[*].PostID
		// 		READ 2: UserTimelineService.WriteUserTimeline() --> usertimeline_cache.*.Get() @ usertimeline_cache.*.Value[*].PostID
		//
		// WARNING: scenarios like those in foreign key read (#3) are problematic and may still occur in other apps
		// CAUSE: problem related to taint propagation:
		// 		- merging old posts (cache key is UserID and never uses exchanged PostID) with new posts (cache values containing exchanged PostID across services) causes the new posts' PostID to "accidently" taint the old posts' PostID
		// 		- the tool then assumes that the PostID was used to read the old posts from their resp. cache
		
		if fread.constraint1 != nil && fread.constraint1.IsMandatory() {
			if !slices.Contains(fread.constraint1.GetRequestsIndexesOnMandatoryFlags(), fread.op1.reqIdx) {
				return true
			}
		} else if fread.constraint2 != nil && fread.constraint2.IsMandatory() {
			if !slices.Contains(fread.constraint2.GetRequestsIndexesOnMandatoryFlags(), fread.op2.reqIdx) {
				return true
			}
		}
	} else {
		return true
	}
	return false
}
