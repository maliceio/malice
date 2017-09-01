package configfile

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

// ConfigFile ~/.docker/config.json file info
type ConfigFile struct {
	Filename string `json:"-"` // Note: for internal use only
}

// SaveToWriter encodes and writes out all the authorization information to
// the given writer
func (configFile *ConfigFile) SaveToWriter(writer io.Writer) error {
	// Encode sensitive data into a new/temp struct
	tmpAuthConfigs := make(map[string]types.AuthConfig, len(configFile.AuthConfigs))
	for k, authConfig := range configFile.AuthConfigs {
		authCopy := authConfig
		// encode and save the authstring, while blanking out the original fields
		authCopy.Auth = encodeAuth(&authCopy)
		authCopy.Username = ""
		authCopy.Password = ""
		authCopy.ServerAddress = ""
		tmpAuthConfigs[k] = authCopy
	}

	saveAuthConfigs := configFile.AuthConfigs
	configFile.AuthConfigs = tmpAuthConfigs
	defer func() { configFile.AuthConfigs = saveAuthConfigs }()

	data, err := json.MarshalIndent(configFile, "", "\t")
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// Save encodes and writes out all the authorization information
func (configFile *ConfigFile) Save() error {
	if configFile.Filename == "" {
		return errors.Errorf("Can't save config with empty filename")
	}

	if err := os.MkdirAll(filepath.Dir(configFile.Filename), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(configFile.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return configFile.SaveToWriter(f)
}
