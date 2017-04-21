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

// Configuration represents the malice runtime configuration.
type Configuration struct {
	Title       string              `toml:"title" json:"title,omitempty"`
	Version     string              `toml:"version" json:"version,omitempty"`
	Author      authorInfo          `toml:"author" json:"author,omitempty"`
	Web         webConfig           `toml:"web" json:"web,omitempty"`
	Email       emailConfig         `toml:"email" json:"email,omitempty"`
	DB          databaseConfig      `toml:"database" json:"db,omitempty"`
	UI          userInterfaceConfig `toml:"ui" json:"ui,omitempty"`
	Environment envConfig           `toml:"environment" json:"environment,omitempty"`
	Docker      dockerConfig        `toml:"docker" json:"docker,omitempty"`
	Logger      loggerConfig        `toml:"logger" json:"logger,omitempty"`
	Proxy       proxyConfig         `toml:"proxy" json:"proxy,omitempty"`
}

type authorInfo struct {
	Name         string `json:"name,omitempty"`
	Organization string `json:"organization,omitempty"`
}

type webConfig struct {
	URL      string `json:"url,omitempty"`
	AdminURL string `toml:"admin_url" json:"admin_url,omitempty"`
}

type userInterfaceConfig struct {
	Name    string `json:"name,omitempty"`
	Image   string `json:"image,omitempty"`
	Server  string `json:"server,omitempty"`
	Ports   []int  `json:"ports,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type databaseConfig struct {
	Name    string `json:"name,omitempty"`
	Image   string `json:"image,omitempty"`
	Server  string `json:"server,omitempty"`
	Ports   []int  `json:"ports,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type emailConfig struct {
	Host     string `json:"host,omitempty"`
	port     int    `json:"port,omitempty"`
	Username string `toml:"user" json:"username,omitempty"`
	Password string `toml:"pass" json:"password,omitempty"`
}

type envConfig struct {
	Run string `json:"run,omitempty"`
}

type dockerConfig struct {
	Name     string        `toml:"machine-name" json:"name,omitempty"`
	EndPoint string        `json:"endpoint,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty"`
	Binds    string        `json:"binds,omitempty"`
	Links    string        `json:"links,omitempty"`
	CPU      int64         `json:"cpu,omitempty"`
	Memory   int64         `json:"memory,omitempty"`
}

type loggerConfig struct {
	FileName   string `json:"filename,omitempty"`
	MaxSize    int    `json:"maxsize,omitempty"`
	MaxAge     int    `json:"maxage,omitempty"`
	MaxBackups int    `json:"maxbackups,omitempty"`
	LocalTime  bool   `json:"localtime,omitempty"`
}

type proxyConfig struct {
	Enable bool   `json:"enable,omitempty"`
	HTTP   string `json:"http,omitempty"`
	HTTPS  string `json:"https,omitempty"`
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
	configPath = path.Join(maldirs.GetConfigDir(), "./config.toml")
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
		er.CheckError(os.MkdirAll(maldirs.GetConfigDir(), 0777))
		// Create the config config in the .malice folder
		er.CheckError(ioutil.WriteFile(configPath, tomlData, 0644))
		log.Debug("Malice config loaded from config/bindata.go")
	}
	er.CheckError(err)

	return
}
