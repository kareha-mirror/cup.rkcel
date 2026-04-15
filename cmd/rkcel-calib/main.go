package main

import (
	"os"
	"path/filepath"

	"tea.kareha.org/cup/rkcel"
	"tea.kareha.org/cup/rkcel/calib"
)

func main() {
	config := rkcel.DefaultConfig()

	var path string
	dir, err := os.UserConfigDir()
	if err == nil {
		path = filepath.Join(dir, "rkcel", "config.yaml")
	}

	if path != "" {
		_, err := os.Stat(path)
		if err == nil { // file exists
			config = rkcel.LoadConfig(path)
		}
	}

	calib.Main(config)

	if path != "" {
		rkcel.SaveConfig(path, config)
	}
}
