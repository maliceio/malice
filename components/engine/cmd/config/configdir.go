package config

import (
	"os"
	"path/filepath"

	"github.com/maliceio/engine/pkg/homedir"
)

var (
	configDir     = os.Getenv("MALICE_CONFIG")
	configFileDir = ".malice"
)

// Dir returns the path to the configuration directory as specified by the MALICE_CONFIG environment variable.
func Dir() string {
	return configDir
}

func init() {
	if configDir == "" {
		configDir = filepath.Join(homedir.Get(), configFileDir)
	}
}
