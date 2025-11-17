package config

type Config struct {
	EnableTransitiveReferences                   bool
	UpdateTransitiveReferencesTriggeredByCurrent bool

	CreateReferencesFromReadReadPair          bool
	CreateReferencesFromReadReadPairAndValKey bool

	PropagateTaintsAcrossQueueOperations bool
}

var Global = &Config{
	EnableTransitiveReferences:                   true,
	UpdateTransitiveReferencesTriggeredByCurrent: true,
	CreateReferencesFromReadReadPair:             false,
	CreateReferencesFromReadReadPairAndValKey:    false,
	PropagateTaintsAcrossQueueOperations:         true,
}
