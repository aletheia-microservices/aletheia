package config

const (
	MIN_T = "t0"
	MAX_T = "t999"
)

type Config struct {
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
}

var Global = &Config{
	// transitive references
	EnableTransitiveReferences:                              false, // do not change
	DeleteOldOnTransitiveReferences:                         false, // tunable
	UpdateTransitiveReferencesTriggeredByCurrent:            false, // do not change
	UpdateTransitiveReferencesTriggeredByCurrentOnMandatory: false, // tunable

	// creation of references on read-read pairs
	CreateReferencesFromReadReadPair:          true, // do not change
	CreateReferencesFromReadReadPairAndValKey: false, // tunable

	// taint propagation
	PropagateTaintsAcrossQueueOperations: true, // do not change
}
