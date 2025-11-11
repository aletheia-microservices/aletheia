package config

type Config struct {
	EnableTransitiveReferences           bool
	CreateReferencesFromReadReadPair     bool
	PropagateTaintsAcrossQueueOperations bool
}

var Global = &Config{
	EnableTransitiveReferences:           true,
	CreateReferencesFromReadReadPair:     false,
	PropagateTaintsAcrossQueueOperations: true,
}
