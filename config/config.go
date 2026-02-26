package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the root configuration structure.
type Config struct {
	Server ServerConfig `yaml:"server"`
	Mongo  MongoConfig  `yaml:"mongo"`
	Kafka  KafkaConfig  `yaml:"kafka"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type MongoConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}

// Load reads the YAML file at the given path and returns a Config.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}
	return &cfg, nil
}
