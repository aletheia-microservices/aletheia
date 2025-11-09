package abstractgraph

import (
	"fmt"

	"analyzer/pkg/common"
	"analyzer/pkg/utils"
)

type AbstractTaint struct {
	fieldpath string // field path in db
	dbcallID  string //i.e., the ID of the abstract edge representing the call
	dbOpType  common.DatabaseOperationType
	primary   bool
	traced    bool
}

func NewAbstractTaint(dbpath string, dbcall string, opType common.DatabaseOperationType, primary bool, traced bool) *AbstractTaint {
	return &AbstractTaint{
		fieldpath: dbpath,
		dbcallID:  dbcall,
		dbOpType:  opType,
		primary:   primary,
		traced:    traced,
	}
}

func (taint *AbstractTaint) Copy() *AbstractTaint {
	return &AbstractTaint{
		fieldpath: taint.fieldpath,
		dbcallID:  taint.dbcallID,
		dbOpType:  taint.dbOpType,
		primary:   taint.primary,
		traced:    taint.traced,
	}
}

func (taint *AbstractTaint) IsRead() bool {
	return taint.dbOpType == common.OP_READ || taint.dbOpType == common.OP_READ_MANY
}

func (taint *AbstractTaint) IsWrite() bool {
	return taint.dbOpType == common.OP_WRITE
}

func (taint *AbstractTaint) IsWriteOrUpdate() bool {
	return taint.dbOpType == common.OP_WRITE || taint.dbOpType == common.OP_UPDATE
}

func (taint *AbstractTaint) IsUpdate() bool {
	return taint.dbOpType == common.OP_UPDATE
}

func (taint *AbstractTaint) IsDelete() bool {
	return taint.dbOpType == common.OP_DELETE
}

func (taint *AbstractTaint) IsPrimary() bool {
	return taint.primary
}

func (taint *AbstractTaint) IsTraced() bool {
	return taint.traced
}

func (taint *AbstractTaint) GetDatabasePath() string {
	return taint.fieldpath
}

func (taint *AbstractTaint) AddSuffixToDatabasePath(suffix string) {
	taint.fieldpath = taint.fieldpath + suffix
}

func (taint *AbstractTaint) GetDatabaseCallID() string {
	return taint.dbcallID
}

func (taint *AbstractTaint) String() string {
	return taint.fieldpath
}

func (taint *AbstractTaint) LongString() string {
	return fmt.Sprintf("{%s, %s, %s, %t}", taint.fieldpath, taint.dbcallID, common.OperationTypeToString(taint.dbOpType), taint.primary)
}

func (taint *AbstractTaint) Similar(other *AbstractTaint) bool {
	fmt.Printf("[ABSTRACT TAINT] [SIMILAR] checking if taints are equal:\n\t%s\n\t%s\n", taint.LongString(), other.LongString())
	return taint.fieldpath == other.fieldpath &&
		taint.dbcallID == other.dbcallID /* &&
		taint.primary == other.primary &&
		taint.dbOpType == other.dbOpType */
}

func (taint *AbstractTaint) Equal(other *AbstractTaint) bool {
	fmt.Printf("[ABSTRACT TAINT] [EQUAL] checking if taints are equal:\n\t%s\n\t%s\n", taint.LongString(), other.LongString())
	return taint.fieldpath == other.fieldpath &&
		taint.dbcallID == other.dbcallID &&
		taint.dbOpType == other.dbOpType &&
		taint.primary == other.primary &&
		taint.traced == other.traced
}

// e.g., 
// - curr dbfield 	= notification
// - other dbfield 	= notification.PostID
func (taint *AbstractTaint) IsUpperTaint(other *AbstractTaint) (bool, string) {
	fmt.Printf("[ABSTRACT TAINT] [SUPER] checking if taint is super path:\n\t%s\n\t%s\n", taint.LongString(), other.LongString())
	if ok, diff := utils.IsUpperPath(taint.fieldpath, other.fieldpath); ok {
		fmt.Printf("got subpath: %s\n", diff)
		return taint.dbcallID == other.dbcallID, diff
	}
	fmt.Printf("[ABSTRACT TAINT] [SUPER] returning false...\n")
	return false, ""
}
