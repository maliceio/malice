package config

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
)

// "github.com/pelletier/go-toml"

// Configuration represents the malice runtime configuration.
type Configuration struct {
	Title       string
	Author      authorInfo
	Web         webConfig
	Email       emailConfig
	DB          databaseConfig      `toml:"database"`
	UI          userInterfaceConfig `toml:ui`
	Environment envConfig
	Docker      dockerConfig
	Logger      loggerConfig
	Proxy       proxyConfig
}

type authorInfo struct {
	Name         string
	Organization string
	Email        string
}

type webConfig struct {
	URL      string
	AdminURL string `toml:"admin_url"`
}

type userInterfaceConfig struct {
	Name    string
	Image   string
	Server  string
	Ports   []int
	Enabled bool
}

type databaseConfig struct {
	Name    string
	Image   string
	Server  string
	Ports   []int
	Timeout int
	Enabled bool
}

type emailConfig struct {
	Host     string
	port     int
	Username string `toml:"user"`
	Password string `toml:"pass"`
}

type envConfig struct {
	Run string
}

type dockerConfig struct {
	Name     string `toml:"machine-name"`
	EndPoint string
	Timeout  time.Duration
	Binds    string
	Links    string
}

type loggerConfig struct {
	FileName   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
}

type proxyConfig struct {
	Enable bool
	HTTP   string
	HTTPS  string
}

// Conf represents the Malice runtime configuration
var Conf Configuration

// Load config.toml into Conf var
// Try to load config from
// - .malice folder       : $HOME/.malice/config.toml
// - binary embedded file : bindata
func Load() {

	var configPath string

	// Check for config config in .malice folder
	configPath = path.Join(maldirs.GetBaseDir(), "./config.toml")
	if _, err := os.Stat(configPath); err == nil {
		_, err := toml.DecodeFile(configPath, &Conf)
		er.CheckError(err)
		log.Debug("Malice config loaded from: ", configPath)
		return
	}

	// Read plugin config out of bindata
	tomlData, err := Asset("config/config.toml")
	if err != nil {
		log.Error(err)
	}
	if _, err = toml.Decode(string(tomlData), &Conf); err == nil {
		// Create .malice folder in the users home directory
		er.CheckError(os.MkdirAll(maldirs.GetBaseDir(), 0777))
		// Create the config config in the .malice folder
		er.CheckError(ioutil.WriteFile(configPath, tomlData, 0644))
		log.Debug("Malice config loaded from config/bindata.go")
	}
	er.CheckError(err)

	return
}
