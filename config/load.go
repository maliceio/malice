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
	DB          databaseConfig `toml:"database"`
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

type databaseConfig struct {
	Path    string
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
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
func Load() {
	// Try to load config from
	// - git repo folder      : MALICE_ROOT/config/config.toml
	// - .malice folder       : $HOME/.malice/config.toml
	// - binary embedded file : bindata

	// Check for config config in repo
	if _, err := os.Stat("./config/config.toml"); err == nil {
		log.Debug("Malice config loaded from ./config/config.toml")

		_, err := toml.DecodeFile("./config/config.toml", &Conf)
		er.CheckError(err)

		return
	}
	// Check for config config in .malice folder
	if _, err := os.Stat(path.Join(maldirs.GetBaseDir(), "./config.toml")); err == nil {
		homeConfigDir := path.Join(maldirs.GetBaseDir(), "./config.toml")
		log.Debug("Malice config loaded from ", homeConfigDir)

		_, err := toml.DecodeFile(homeConfigDir, &Conf)
		er.CheckError(err)

		return
	}
	// Read plugin config out of bindata
	tomlData, err := Asset("config/config.toml")
	if err != nil {
		log.Error(err)
	}

	if _, err := toml.Decode(string(tomlData), &Conf); err == nil {
		log.Debug("Malice config loaded from config/bindata.go")
		configPath := path.Join(maldirs.GetBaseDir(), "./config.toml")
		er.CheckError(ioutil.WriteFile(configPath, tomlData, 0644))
	}
	er.CheckError(err)

	return
}
