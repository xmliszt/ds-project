package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ConfigServer `yaml:"Server"`
	ConfigNode `yaml:"Node"`
}

type ConfigServer struct {
	Host string `yaml:"Host"`
	Port int `yaml:"Port"`
}

type ConfigNode struct {
	Number int `yaml:"Number"`
	HeartbeatInterval int `yaml:"HeartbeatInterval"`
}

var configPath = "./config.yaml"

// LoadConfig loads the config from YAML file and return the config object
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	d := yaml.NewDecoder(file)

	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsEqual compares nested struct
func (config *Config) IsEqual(xConfig *Config) bool {
	return config.ConfigServer == xConfig.ConfigServer && config.ConfigNode == xConfig.ConfigNode 
}

