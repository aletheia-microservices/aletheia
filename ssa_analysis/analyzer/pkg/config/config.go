package config

type Config struct {
	EnableTransitiveReferences                bool
	CreateReferencesFromReadReadPair          bool
	CreateReferencesFromReadReadPairAndValKey bool
	PropagateTaintsAcrossQueueOperations      bool
}

var Global = &Config{
	EnableTransitiveReferences:                true,
	CreateReferencesFromReadReadPair:          false,
	CreateReferencesFromReadReadPairAndValKey: false,
	PropagateTaintsAcrossQueueOperations:      true,
}
