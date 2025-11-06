package config

type Config struct {
	EnableTransitiveReferences bool
}

var Global = &Config{
	EnableTransitiveReferences: true,
}
