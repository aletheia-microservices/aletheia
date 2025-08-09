package abstractgraph

import (
	"fmt"
	"strings"

	"analyzer/pkg/common"
)

type AbstractTaint struct {
	dbpath   string // field path in db
	dbcallID string //i.e., the ID of the abstract edge representing the call
	dbOpType common.DatabaseOperationType
	primary  bool
	traced   bool
}

func NewAbstractTaint(dbpath string, dbcall string, opType common.DatabaseOperationType, primary bool, traced bool) *AbstractTaint {
	return &AbstractTaint{
		dbpath:   dbpath,
		dbcallID: dbcall,
		dbOpType: opType,
		primary:  primary,
		traced:   traced,
	}
}

func (taint *AbstractTaint) IsRead() bool {
	return taint.dbOpType == common.OP_READ
}

func (taint *AbstractTaint) IsWrite() bool {
	return taint.dbOpType == common.OP_WRITE
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
	return taint.dbpath
}

func (taint *AbstractTaint) GetDatabaseCallID() string {
	return taint.dbcallID
}

func (taint *AbstractTaint) String() string {
	return taint.dbpath
}

func (taint *AbstractTaint) LongString() string {
	return fmt.Sprintf("{%s, %s, %s, %t}", taint.dbpath, taint.dbcallID, common.OperationTypeToString(taint.dbOpType), taint.primary)
}

func (taint *AbstractTaint) Equals(other *AbstractTaint) bool {
	fmt.Printf("[ABSTRACT TAINT] [EQUAL] checking if taints are equal:\n\t%s\n\t%s\n", taint.LongString(), other.LongString())
	return taint.dbpath == other.dbpath &&
		taint.dbcallID == other.dbcallID /* &&
		taint.primary == other.primary &&
		taint.dbOpType == other.dbOpType */
}

// taint.dbfield: notification
// other.dbfield: notification.PostID
func (taint *AbstractTaint) IsUpperPath(other *AbstractTaint) (bool, string) {
	fmt.Printf("[ABSTRACT TAINT] [SUPER] checking if taint is super path:\n\t%s\n\t%s\n", taint.LongString(), other.LongString())
	if taint.dbpath != other.dbpath && strings.HasPrefix(other.dbpath, taint.dbpath) {
		var subpath string
		_, subpath, _ = strings.Cut(other.dbpath, taint.dbpath)
		fmt.Printf("got subpath: %s\n", subpath)
		return taint.dbcallID == other.dbcallID, subpath /* &&
			taint.primary == other.primary &&
			taint.dbOpType == other.dbOpType, subpath */
	}
	return false, ""
}
