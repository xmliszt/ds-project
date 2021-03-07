package config

type GlobalConfig struct {
	NodeNumber int
}

var DefaultConfig = GlobalConfig{
	NodeNumber: 5,
}

func LoadConfig() *GlobalConfig {
	cfg := &GlobalConfig{}
	*cfg = DefaultConfig
	return cfg
}