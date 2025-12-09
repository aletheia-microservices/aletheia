package foreignkeyconcurrency

import (
	"analyzer/pkg/app/backends"
)

type DangerousDelete struct {
	delete           *DeleteOperation
	concurrentWrites []*ConcurrentWrite
}

func (dd *DangerousDelete) CallString() string {
	return dd.delete.call.String()
}

type ConcurrentWrite struct {
	write          *WriteOperation
	affectedFields []*backends.Field
	database       *backends.Database
	schema         *backends.Schema
}

func (cw *ConcurrentWrite) CallString() string {
	return cw.write.call.String()
}

func (cw *ConcurrentWrite) EntryString() string {
	return cw.write.request.entry.String()
}

func (detector *ForeignKeyConcurrencyDetector) checkInconsistencies() {
	// EVAL: logrus.Tracef("[FOREIGN KEY CONCURRENCY | CHECKER] checking inconsistencies\n")
	for _, request := range detector.requests {
		for _, delete := range request.getAllDeleteOperations() {
			// EVAL: logrus.Tracef("\t[FOREIGN KEY CONCURRENCY | CHECKER] delete = %s\n", delete.call.String())
			var concurrentWrites map[*WriteOperation][]*backends.Field

			for _, otherRequest := range detector.requests {
				if otherRequest.idx == request.idx {
					continue
				}
				for _, otherWrite := range otherRequest.getAllWriteOperations() {
					// EVAL: logrus.Tracef("\t[FOREIGN KEY CONCURRENCY | CHECKER] other_write={%s}, entry={%s}\n", otherWrite.call.String(), otherWrite.request.entry.String())
					for _, otherField := range otherWrite.fields {
						// EVAL: logrus.Tracef("\t\t[FOREIGN KEY CONCURRENCY | CHECKER] other field = %s\n", otherField.String())
						for _, deletedField := range delete.schema.GetAllFieldsLst() {
							// EVAL: logrus.Tracef("\t\t[FOREIGN KEY CONCURRENCY | CHECKER] deleted field = %s\n", deletedField.String())
							if otherField.HasConstraintForeignKeyToField(deletedField) {
								if concurrentWrites == nil {
									concurrentWrites = make(map[*WriteOperation][]*backends.Field)
								}
								concurrentWrites[otherWrite] = append(concurrentWrites[otherWrite], otherField)
								// EVAL: logrus.Tracef("\t\t\t[FOREIGN KEY CONCURRENCY | CHECKER] OK!\n")
							}
						}
					}
				}
			}

			if concurrentWrites != nil {
				dangerousDelete := &DangerousDelete{
					delete: delete,
				}
				for concurrentWrite, affectedFields := range concurrentWrites {
					concurrentWrite := &ConcurrentWrite{
						write:          concurrentWrite,
						affectedFields: affectedFields,
					}
					concurrentWrite.database = affectedFields[0].GetDatabase()
					concurrentWrite.schema = affectedFields[0].GetSchema()
					dangerousDelete.concurrentWrites = append(dangerousDelete.concurrentWrites, concurrentWrite)
				}
				detector.addDangerousDelete(request, dangerousDelete)
			}
		}
	}
}
