package config

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Config config
type Config struct {
	CommonConfig
}

// CommonConfig common config
type CommonConfig struct {
	Debug     bool     `json:"debug,omitempty"`
	Hosts     []string `json:"hosts,omitempty"`
	LogLevel  string   `json:"log-level,omitempty"`
	TLS       bool     `json:"tls,omitempty"`
	TLSVerify bool     `json:"tlsverify,omitempty"`
}

// New returns a new fully initialized Config struct
func New() *Config {
	config := Config{}

	return &config
}

// Reload reads the configuration in the host and reloads the daemon and server.
func Reload(newConfig *Config, flags *pflag.FlagSet, reload func(*Config)) error {
	logrus.Infof("Got signal to reload configuration, reloading from: %s", newConfig)

	if err := Validate(newConfig); err != nil {
		return fmt.Errorf("file configuration validation failed (%v)", err)
	}

	reload(newConfig)
	return nil
}

// Validate validates some specific configs.
func Validate(config *Config) error {
	return nil
}
