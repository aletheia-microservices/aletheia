package config

type Config struct {
	EnableTransitiveReferences bool
	CreateReferencesFromReadReadPair bool
}

var Global = &Config{
	EnableTransitiveReferences: true,
	CreateReferencesFromReadReadPair: false,
}
