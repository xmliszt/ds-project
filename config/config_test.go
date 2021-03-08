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
		HeartbeatInterval: 5,
	},
	ConfigTimeout: ConfigTimeout{
		HeartBeatTimeout: 15,
		NodeCreationTimeout: 60,
	},
}

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Error(err)
	}
	if !(config.ConfigServer == expectedConfig.ConfigServer && config.ConfigNode == expectedConfig.ConfigNode) {
		t.Errorf("Expected: %v, instead received: %v", expectedConfig, config)
	}
}