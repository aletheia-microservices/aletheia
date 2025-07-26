package abstractcallgraph

import "fmt"

type AbstractTaint struct {
	dbfield  string
	dbcallID string //i.e., the ID of the abstract edge representing the call
}

func NewAbstractTaint(dbfield string, dbcall string) *AbstractTaint {
	return &AbstractTaint{
		dbfield:  dbfield,
		dbcallID: dbcall,
	}
}

func (taint *AbstractTaint) GetDbField() string {
	return taint.dbfield
}

func (taint *AbstractTaint) GetDbCall() string {
	return taint.dbcallID
}

func (taint *AbstractTaint) String() string {
	if taint.dbcallID == "" {
		return fmt.Sprintf("%s (call at <nil>)", taint.dbfield)
	}
	return fmt.Sprintf("%s (call at %s)", taint.dbfield, taint.dbcallID)
}
