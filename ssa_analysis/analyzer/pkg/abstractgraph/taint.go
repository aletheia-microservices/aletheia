package abstractgraph

import "fmt"

type AbstractTaint struct {
	dbfield  string // field path in db
	dbcallID string //i.e., the ID of the abstract edge representing the call
	primary  bool
	write    bool
}

func NewAbstractTaint(dbfield string, dbcall string, primary bool, write bool) *AbstractTaint {
	return &AbstractTaint{
		dbfield:  dbfield,
		dbcallID: dbcall,
		primary:  primary,
		write:    write,
	}
}

func (taint *AbstractTaint) IsWrite() bool {
	return taint.write
}

func (taint *AbstractTaint) IsPrimary() bool {
	return taint.primary
}

func (taint *AbstractTaint) GetField() string {
	return taint.dbfield
}

func (taint *AbstractTaint) GetCallID() string {
	return taint.dbcallID
}

func (taint *AbstractTaint) String() string {
	return taint.dbfield
}

func (taint *AbstractTaint) LongString() string {
	return fmt.Sprintf("{%s, %s, %t, %t}", taint.dbfield, taint.dbcallID, taint.primary, taint.write)
}

func (taint *AbstractTaint) Equals(other *AbstractTaint) bool {
	fmt.Printf("checking if taints are equal:\n\t%s\n\t%s\n", taint.LongString(), other.LongString())
	return taint.dbfield == other.dbfield &&
		taint.dbcallID == other.dbcallID &&
		taint.primary == other.primary &&
		taint.write == other.write
}
