package config

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/maliceio/malice/libmalice/maldirs"
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
}

type loggerConfig struct {
	FileName   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Conf represents the Malice runtime configuration
var Conf Configuration

func init() {
	// Get the config file
	configPath := "./config.toml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = path.Join(maldirs.GetBaseDir(), configPath)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			assert(err)
		}
	}

	_, err := toml.DecodeFile(configPath, &Conf)
	assert(err)
	// fmt.Println(Conf)
}
