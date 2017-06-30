package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/maliceio/malice/utils"
)

// Configuration represents the malice runtime configuration.
type Configuration struct {
	Title       string              `toml:"title"`
	Version     string              `toml:"version"`
	Author      authorInfo          `toml:"author"`
	Web         webConfig           `toml:"web"`
	Email       emailConfig         `toml:"email"`
	DB          databaseConfig      `toml:"database"`
	UI          userInterfaceConfig `toml:"ui"`
	Environment envConfig           `toml:"environment"`
	Docker      dockerConfig        `toml:"docker"`
	Logger      loggerConfig        `toml:"logger"`
	Proxy       proxyConfig         `toml:"proxy"`
}

type authorInfo struct {
	Name         string `toml:"name"`
	Organization string `toml:"organization"`
}

type webConfig struct {
	URL      string `toml:"url"`
	AdminURL string `toml:"admin_url"`
}

type userInterfaceConfig struct {
	Name    string `toml:"name"`
	Image   string `toml:"image"`
	Server  string `toml:"server"`
	Ports   []int  `toml:"ports"`
	Enabled bool   `toml:"enabled"`
}

type databaseConfig struct {
	Name    string `toml:"name"`
	Image   string `toml:"image"`
	Server  string `toml:"server"`
	Ports   []int  `toml:"ports"`
	Timeout int    `toml:"timeout"`
	Enabled bool   `toml:"enabled"`
}

type emailConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"user"`
	Password string `toml:"pass"`
}

type envConfig struct {
	Run string `toml:"run"`
}

type dockerConfig struct {
	Name     string        `toml:"machine-name"`
	EndPoint string        `toml:"endpoint"`
	Timeout  time.Duration `toml:"timeout"`
	Binds    string        `toml:"binds"`
	Links    string        `toml:"links"`
	CPU      int64         `toml:"cpu"`
	Memory   int64         `toml:"memory"`
}

type loggerConfig struct {
	FileName   string `toml:"filename"`
	MaxSize    int    `toml:"maxsize"`
	MaxAge     int    `toml:"maxage"`
	MaxBackups int    `toml:"maxbackups"`
	LocalTime  bool   `toml:"localtime"`
}

type proxyConfig struct {
	Enable bool   `toml:"enable"`
	HTTP   string `toml:"http"`
	HTTPS  string `toml:"https"`
}

// Conf represents the Malice runtime configuration
var Conf Configuration

// UpdateConfig will update the config on disk with the one embedded in malice
func UpdateConfig() error {
	configPath := path.Join(maldirs.GetConfigDir(), "./config.toml")
	configBackupPath := path.Join(maldirs.GetConfigDir(), "./config.toml.backup")
	er.CheckError(os.Rename(configPath, configBackupPath))
	// Read plugin config out of bindata
	tomlData, err := Asset("config/config.toml")
	if err != nil {
		log.Error(err)
	}
	if _, err = toml.Decode(string(tomlData), &Conf); err == nil {
		// Update the config config in the .malice folder
		er.CheckError(ioutil.WriteFile(configPath, tomlData, 0644))
		log.Debug("Malice config loaded from config/bindata.go")
	}
	return err
}

func LoadFromToml(configPath, version string) {
	_, err := toml.DecodeFile(configPath, &Conf)
	if err != nil {
		// try the config embedded in malice instead
		log.Debug("Malice config loaded from embedded binary data")
		loadFromBinary(configPath)
	}
	log.Debug("Malice config loaded from: ", configPath)
	log.Debugf("config.toml version: %s, malice version: %s", Conf.Version, version)
	if version != "dev" && !strings.Contains(Conf.Version, version) {
		// Prompt user to update malice config.toml?
		log.Infof("Newer version of malice config.toml available: %s, you currently have %s", version, Conf.Version)
		fmt.Println("Would you like to update now? (yes/no)")
		if utils.AskForConfirmation() {
			log.Debug("Updating config: ", configPath)
			er.CheckError(UpdateConfig())
		}
	}
}

func loadFromBinary(configPath string) {
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
}

// Load config.toml into Conf var
// Try to load config from
// - .malice folder       : $HOME/.malice/config.toml
// - binary embedded file : bindata
func Load(version string) {
	// Check for config config in .malice folder
	configPath := path.Join(maldirs.GetConfigDir(), "./config.toml")
	if _, err := os.Stat(configPath); err == nil {
		LoadFromToml(configPath, version)
		return
	}
	loadFromBinary(configPath)
	return
}
