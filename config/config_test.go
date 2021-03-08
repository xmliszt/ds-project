package config

import (
	"os"
	"path"
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
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	rootPath, _ := path.Split(cwd)
	configPath := path.Join(rootPath, "config.yaml")
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("Error opening config! ERROR: %v | Config File Path: %s", err, configPath)
	}
	if !(config.ConfigServer == expectedConfig.ConfigServer && config.ConfigNode == expectedConfig.ConfigNode) {
		t.Errorf("Expected: %v, instead received: %v", expectedConfig, config)
	}
}