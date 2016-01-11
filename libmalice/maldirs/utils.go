package maldirs

import (
	"os"
	"path/filepath"

	"github.com/maliceio/malice/libmalice/malutils"
)

// NOTE: export MALICE_STORAGE_PATH=$(pwd)
var (
	BaseDir = os.Getenv("MALICE_STORAGE_PATH")
)

func GetBaseDir() string {
	if BaseDir == "" {
		BaseDir = filepath.Join(malutils.GetHomeDir(), ".malice")
	}
	return BaseDir
}

func GetSampledsDir() string {
	return filepath.Join(GetBaseDir(), "samples")
}

func GetPluginsDir() string {
	return filepath.Join(GetBaseDir(), "plugins")
}

func GetConfigDir() string {
	return filepath.Join(GetBaseDir(), "config")
}
