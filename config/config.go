package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// "github.com/pelletier/go-toml"

// Configuration represents the malice runtime configuration.
type Configuration struct {
	Title       string
	Author      authorInfo
	Web         webConfig
	Email       emailConfig
	DB          databaseConfig    `toml:"database"`
	Plugins     map[string]Plugin `toml:"plugin"`
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

// Plugin represents a single plugin setting.
type Plugin struct {
	Enabled     bool
	Category    string
	Description string
	Image       string
	Mime        string
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
	_, err := toml.DecodeFile("./config.toml", &Conf)
	assert(err)
	// fmt.Println(Conf)
}
