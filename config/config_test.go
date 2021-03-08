package config

import (
	"testing"
)

var expectedConfig = &Config{
	ConfigServer: ConfigServer{
		Host: "localhost",
		Port: 8080,
	},
	ConfigNode: ConfigNode{
		Number: 5,
	},
}

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Error(err)
	}
	if !config.IsEqual(expectedConfig) {
		t.Errorf("Expected: %v, instead received: %v", expectedConfig, config)
	}
}