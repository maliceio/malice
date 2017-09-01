package config

import (
	"os"
	"path/filepath"

	"github.com/maliceio/engine/cli/config/configfile"
	"github.com/maliceio/engine/pkg/homedir"
)

const (
	// ConfigFileName is the name of config file
	ConfigFileName = "config.json"
	configFileDir  = ".malice"
)

var (
	configDir = os.Getenv("MALICE_CONFIG")
)

func init() {
	if configDir == "" {
		configDir = filepath.Join(homedir.Get(), configFileDir)
	}
}

// Dir returns the directory the configuration file is stored in
func Dir() string {
	return configDir
}

// SetDir sets the directory the configuration file is stored in
func SetDir(dir string) {
	configDir = dir
}

// NewConfigFile initializes an empty configuration file for the given filename 'fn'
func NewConfigFile(fn string) *configfile.ConfigFile {
	return &configfile.ConfigFile{
		// AuthConfigs: make(map[string]types.AuthConfig),
		// HTTPHeaders: make(map[string]string),
		Filename: fn,
	}
}
