package rkcel

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CellWidth  int `yaml:"cell-width"`
	CellHeight int `yaml:"cell-height"`
}

func DefaultConfig() *Config {
	return &Config{
		CellWidth:  6,
		CellHeight: 13,
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveConfig(path string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(path), 0755)

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
