package datastores

type DatabaseInstance interface {
	String() string
	LongString() string
	IsQueue() bool
	GetTypeName() string
	GetTypeLongName() string
	GetName() string
	Equals() bool
	GetDatastore() *Datastore
}
