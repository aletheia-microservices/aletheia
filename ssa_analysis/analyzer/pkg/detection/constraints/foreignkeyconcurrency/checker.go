package foreignkeyconcurrency

import (
	"fmt"

	"analyzer/pkg/app/backends"
)

type DangerousDelete struct {
	delete           *DeleteOperation
	concurrentWrites []*ConcurrentWrite
}

type ConcurrentWrite struct {
	write          *WriteOperation
	affectedFields []*backends.Field
	database       *backends.Database
	schema         *backends.Schema
}

func (detector *ForeignKeyConcurrencyDetector) checkInconsistencies() {
	fmt.Printf("[FOREIGN KEY CONCURRENCY | CHECKER] checking inconsistencies\n")
	for _, request := range detector.requests {
		for _, delete := range request.getAllDeleteOperations() {
			fmt.Printf("\t[FOREIGN KEY CONCURRENCY | CHECKER] delete = %s\n", delete.call.String())
			var concurrentWrites map[*WriteOperation][]*backends.Field

			for _, otherRequest := range detector.requests {
				if otherRequest.idx == request.idx {
					continue
				}
				for _, otherWrite := range otherRequest.getAllWriteOperations() {
					fmt.Printf("\t[FOREIGN KEY CONCURRENCY | CHECKER] other write = %s\n", otherWrite.call.String())
					for _, otherField := range otherWrite.fields {
						fmt.Printf("\t\t[FOREIGN KEY CONCURRENCY | CHECKER] other field = %s\n", otherField.String())

						for _, deletedField := range delete.schema.GetAllFieldsLst() {
							fmt.Printf("\t\t[FOREIGN KEY CONCURRENCY | CHECKER] deleted field = %s\n", deletedField.String())
							if otherField.HasConstraintForeignKeyNonMandatoryToField(deletedField) {
								if concurrentWrites == nil {
									concurrentWrites = make(map[*WriteOperation][]*backends.Field)
								}
								concurrentWrites[otherWrite] = append(concurrentWrites[otherWrite], otherField)
								fmt.Printf("\t\t\t[FOREIGN KEY CONCURRENCY | CHECKER] OK!\n")
							}
						}
					}
				}
			}

			if concurrentWrites != nil {
				dangerousDelete := &DangerousDelete{
					delete: delete,
				}
				for write, affectedFields := range concurrentWrites {
					concurrentWrite := &ConcurrentWrite{
						write:          write,
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
