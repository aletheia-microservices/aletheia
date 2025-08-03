package foreignkeycoordination

type ForeignRead struct {
	read1 *ReadOperation // read1 has the secondary taint referencing read2
	read2 *ReadOperation
}

func (detector *ForeignKeyCoordinationDetector) checkInconsistency(request *Request, read *ReadOperation) {
	for _, arg := range read.arguments {
		for _, secondaryTaint := range arg.GetSecondaryTaintsFlatList() {
			if secondaryTaint.GetCallID() != read.GetCallID() {
				otherRead := request.FindOperationByCallID(secondaryTaint.GetCallID())
				if otherRead != nil {
					foreignRead := &ForeignRead{read1: read, read2: otherRead}
					detector.addInconsistency(request, foreignRead)
				}
			}
		}
	}
}
