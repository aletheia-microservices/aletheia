package foreignkeycoordination

type ForeignRead struct {
	op1 *ReadOperation // read1 has the secondary taint referencing read2
	op2 *ReadOperation
}

func (detector *ForeignKeyCoordinationDetector) checkInconsistency(request *Request, currOp *ReadOperation) {
	// same logic as in foreignkeycascade but here we verify if secondaryTaint.IsDelete()
	for _, arg := range currOp.arguments {
		for _, secondaryTaint := range arg.GetSecondaryTaintsFlatList() {
			if secondaryTaint.GetCallID() != currOp.GetCallID() && secondaryTaint.IsRead() {
				otherOp := request.FindOperationByCallID(secondaryTaint.GetCallID())
				if otherOp != nil {
					foreignRead := &ForeignRead{op1: currOp, op2: otherOp}
					detector.addForeignRead(request, foreignRead)
				}
			}
		}
	}
}
