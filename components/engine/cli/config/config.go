package config

import (
	"os"
	"path/filepath"

	"github.com/maliceio/engine/cli/config/configfile"
	"github.com/maliceio/engine/pkg/homedir"
	"github.com/pkg/errors"
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

// Load reads the configuration files in the given directory, and sets up
// the auth config information and returns values.
// FIXME: use the internal golang config parser
func Load(configDir string) (*configfile.ConfigFile, error) {
	if configDir == "" {
		configDir = Dir()
	}

	filename := filepath.Join(configDir, ConfigFileName)
	configFile := configfile.New(filename)

	// Try happy path first - latest config file
	if _, err := os.Stat(filename); err == nil {
		file, err := os.Open(filename)
		if err != nil {
			return configFile, errors.Errorf("%s - %v", filename, err)
		}
		defer file.Close()
		err = configFile.LoadFromReader(file)
		if err != nil {
			err = errors.Errorf("%s - %v", filename, err)
		}
		return configFile, err
	} else if !os.IsNotExist(err) {
		// if file is there but we can't stat it for any reason other
		// than it doesn't exist then stop
		return configFile, errors.Errorf("%s - %v", filename, err)
	}
	return configFile, nil
}
