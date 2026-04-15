package rkcel

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CellWidth  int  `yaml:"cell-width"`
	CellHeight int  `yaml:"cell-height"`
	UseBottom  bool `yaml:"use-bottom"`
}

func DefaultConfig() *Config {
	return &Config{
		CellWidth:  8,
		CellHeight: 16,
		UseBottom:  false,
	}
}

func LoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return &cfg
}

func SaveConfig(path string, config *Config) {
	data, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("failed to serialize config: %v", err)
	}

	os.MkdirAll(filepath.Dir(path), 0755)

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatalf("failed to save config: %v", err)
	}
}
