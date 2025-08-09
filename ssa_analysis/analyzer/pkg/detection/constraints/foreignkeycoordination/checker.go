package foreignkeycoordination

import "fmt"

type ForeignRead struct {
	op1 *ReadOperation // read1 has the secondary taint referencing read2
	op2 *ReadOperation
}

func (detector *ForeignKeyCoordinationDetector) checkInconsistency(request *Request, currOp *ReadOperation) {
	// same logic as in foreignkeycascade but here we verify if secondaryTaint.IsDelete()
	for _, arg := range currOp.arguments {
		fmt.Printf("[FOREIGN KEY CHECKER] arg (%s) on op: %v\n", arg.String(), currOp.call.String())
		fmt.Printf("\t[FOREIGN KEY CHECKER] arg (%s) w/ all taints: %v\n", arg.String(), arg.GetAllTaints())
		for _, secondaryTaint := range arg.GetSecondaryTaintsFlatList() {
			fmt.Printf("\t[FOREIGN KEY CHECKER] arg (%s) w/ secondary taints: %v\n", arg.String(), secondaryTaint.LongString())
			if secondaryTaint.GetDatabaseCallID() != currOp.GetCallID() && secondaryTaint.IsRead() {
				otherOp := request.FindOperationByCallID(secondaryTaint.GetDatabaseCallID())
				if otherOp != nil && !detector.hasForeignRead(request, currOp, otherOp) {
					foreignRead := &ForeignRead{op1: currOp, op2: otherOp}
					detector.addForeignRead(request, foreignRead)
				}
			}
		}
	}
}
