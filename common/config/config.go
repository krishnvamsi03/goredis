package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Loginfo       Loginfo       `yaml:"log"`
		ServerOptions ServerOptions `yaml:"srvoptions"`
	}

	Loginfo struct {
		Level string `yaml:"level"`
	}

	ServerOptions struct {
		Port string `yaml:"port"`
	}
)

var defaultConfigPath string = "./config/config.yaml"

func LoadConfig(configPath string) (*Config, error) {

	if len(configPath) == 0 {
		configPath = defaultConfigPath
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
