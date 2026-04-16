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

	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func UserConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "rkcel", "config.yaml")
	return path, nil
}

func LoadUserConfig() (*Config, error) {
	config := DefaultConfig()

	path, err := UserConfigPath()
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(path)
	if err == nil { // file exists
		config, err = LoadConfig(path)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func SaveUserConfig(config *Config) error {
	path, err := UserConfigPath()
	if err != nil {
		return err
	}
	return SaveConfig(path, config)
}
