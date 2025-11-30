package config

const (
	MIN_T = "t0"
	MAX_T = "t999"
)

type Config struct {
	// pattern detection
	RestrictivePrimaryKeyCoordinationAnalysis bool
	RestrictiveForeignKeyCoordinationAnalysis bool

	// transitive references
	EnableTransitiveReferences                              bool
	DeleteOldOnTransitiveReferences                         bool
	UpdateTransitiveReferencesTriggeredByCurrent            bool
	UpdateTransitiveReferencesTriggeredByCurrentOnMandatory bool

	// creation of references on read-read pairs
	CreateReferencesFromReadReadPair          bool
	CreateReferencesFromReadReadPairAndValKey bool

	// taint propagation
	PropagateTaintsAcrossQueueOperations bool

	// database schema configuration
	MakeIndexesAsPrimaryKeysForNoSQLDatabases bool

	ForeignKeyConcurrencyDetectorIncludeOnUpdates bool
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

	MakeIndexesAsPrimaryKeysForNoSQLDatabases: true, // do not change

	ForeignKeyConcurrencyDetectorIncludeOnUpdates: false, // tunable
}
