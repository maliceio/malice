package config

import (
	"io/ioutil"
	"os"
	"time"

	"path"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/data"
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
	// Get the config file
	configPath := "./data/config.toml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// er.CheckErrorNoStackWithMessage(err, "NOT FOUND")
		configPath = path.Join(maldirs.GetBaseDir(), "./config.toml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configData, err := data.Asset("data/config.toml")
			er.CheckError(err)
			er.CheckError(ioutil.WriteFile(configPath, configData, 0644))
		}
	}
	log.Debug("Malice Config: ", configPath)
	_, err := toml.DecodeFile(configPath, &Conf)
	er.CheckError(err)
	// fmt.Println(Conf)
}
