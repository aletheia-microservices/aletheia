package abstractcallgraph

type AbstractTaint struct {
	dbfield  string
	dbcallID string //i.e., the ID of the abstract edge representing the call
	primary  bool
}

func NewAbstractTaint(dbfield string, dbcall string, primary bool) *AbstractTaint {
	return &AbstractTaint{
		dbfield:  dbfield,
		dbcallID: dbcall,
		primary:  primary,
	}
}

func (taint *AbstractTaint) IsPrimary() bool {
	return taint.primary
}

func (taint *AbstractTaint) GetDbField() string {
	return taint.dbfield
}

func (taint *AbstractTaint) GetDbCall() string {
	return taint.dbcallID
}

func (taint *AbstractTaint) String() string {
	return taint.dbfield
}

func (taint *AbstractTaint) Equals(other *AbstractTaint) bool {
	return taint.dbfield == other.dbfield && taint.dbcallID == other.dbcallID && taint.primary == other.primary
}
