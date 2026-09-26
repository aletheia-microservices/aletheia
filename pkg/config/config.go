// Package config holds the global analysis settings (see Config)
package config

const (
	MIN_T = "t0"
	MAX_T = "t999"
)

// Config holds global analysis settings. In Global, settings marked "tunable" can be changed
// to adjust the analysis, while settings marked "do not change" are required for the analysis
// to work as described in the paper
type Config struct {
	// --- pattern detection
	// EI-1 only reports primary key reads linked by a mandatory foreign key created in a different request
	RestrictivePrimaryKeyCoordinationAnalysis bool
	// RI-3 only reports mandatory foreign keys created in a different request than the reads
	RestrictiveForeignKeyCoordinationAnalysis bool

	// --- transitive references
	// infer transitive foreign keys (X references Y and Y references Z, so X references Z)
	EnableTransitiveReferences bool
	// remove X references Y once it has been extended into X references Z
	DeleteOldOnTransitiveReferences bool
	// when a new foreign key Y references Z is inferred, extend existing foreign keys X references Y into X references Z
	UpdateTransitiveReferencesTriggeredByCurrent bool
	// also extend X references Y when it is mandatory
	UpdateTransitiveReferencesTriggeredByCurrentOnMandatory bool

	// --- creation of references on read-read pairs
	// infer foreign keys from (read_key, read_key) pairs
	CreateReferencesFromReadReadPair bool
	// also infer foreign keys from (read_val, read_key) pairs
	CreateReferencesFromReadReadPairAndValKey bool

	// --- taint propagation
	// follow values from a queue push to the matching pop
	PropagateTaintsAcrossQueueOperations bool

	// --- database schema configuration
	// treat the uniqueItems of NoSQL collections as primary keys, not only as unique
	MakeIndexesAsPrimaryKeysForNoSQLDatabases bool
	// RI-2 also treats NoSQL Upsert operations as writes
	ForeignKeyConcurrencyDetectorIncludeOnUpdates bool
	// build the schema in two passes: write pairs first, then read-read pairs
	DualPassSchemaBuilding bool
}

var Global = &Config{
	// pattern detection
	RestrictivePrimaryKeyCoordinationAnalysis: true,
	RestrictiveForeignKeyCoordinationAnalysis: true,

	// transitive references
	EnableTransitiveReferences:                              true,  // do not change
	DeleteOldOnTransitiveReferences:                         true,  // tunable
	UpdateTransitiveReferencesTriggeredByCurrent:            true,  // do not change
	UpdateTransitiveReferencesTriggeredByCurrentOnMandatory: false, // tunable

	// creation of references on read-read pairs
	CreateReferencesFromReadReadPair:          true,  // do not change
	CreateReferencesFromReadReadPairAndValKey: false, // tunable

	// taint propagation
	PropagateTaintsAcrossQueueOperations: true, // do not change

	MakeIndexesAsPrimaryKeysForNoSQLDatabases: false, // do not change

	ForeignKeyConcurrencyDetectorIncludeOnUpdates: false, // tunable

	DualPassSchemaBuilding: false, // do not change
}
