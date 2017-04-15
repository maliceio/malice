package maldirs

import (
	"os"
	"path/filepath"

	"github.com/maliceio/malice/malice/malutils"
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

func GetLogsDir() string {
	return filepath.Join(GetBaseDir(), "logs")
}

func MakeDirs() {
	// Make .malice directory if it doesn't exist
	if _, err := os.Stat(GetSampledsDir()); os.IsNotExist(err) {
		os.MkdirAll(GetSampledsDir(), 0777)
	}
	if _, err := os.Stat(GetPluginsDir()); os.IsNotExist(err) {
		os.MkdirAll(GetPluginsDir(), 0777)
	}
	if _, err := os.Stat(GetConfigDir()); os.IsNotExist(err) {
		os.MkdirAll(GetConfigDir(), 0777)
	}
	if _, err := os.Stat(GetLogsDir()); os.IsNotExist(err) {
		os.MkdirAll(GetLogsDir(), 0777)
	}
}
