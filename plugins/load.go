package plugins

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
)

// Plugin represents a single plugin setting.
type Plugin struct {
	Name        string   `toml:"name"`
	Enabled     bool     `toml:"enabled"`
	Category    string   `toml:"category"`
	Description string   `toml:"description"`
	Image       string   `toml:"image"`
	Repository  string   `toml:"repository"`
	Build       bool     `toml:"build"`
	APIKey      string   `toml:"apikey"`
	Mime        string   `toml:"mime"`
	HashTypes   []string `toml:"hashtypes"`
	Cmd         string   `toml:"cmd"`
	Env         []string `toml:"env"`
	Installed   bool
}

// Configuration represents the malice runtime plugins.
type Configuration struct {
	Plugins []Plugin `toml:"plugin"`
}

// Plugs represents the Malice runtime configuration
var Plugs Configuration

// Load plugins.toml into Plug var
// Try to load plugins from
// - .malice folder       : $HOME/.malice/plugins.toml
// - binary embedded file : bindata
func Load() {

	var configPath string

	// Check for plugins config in .malice folder
	configPath = path.Join(maldirs.GetPluginsDir(), "./plugins.toml")
	if _, err := os.Stat(configPath); err == nil {
		_, err := toml.DecodeFile(configPath, &Plugs)
		er.CheckError(err)
		log.Debug("Malice plugins loaded from: ", configPath)
		return
	}

	// Read plugin config out of bindata
	tomlData, err := Asset("plugins/plugins.toml")
	if err != nil {
		log.Error(err)
	}
	if _, err = toml.Decode(string(tomlData), &Plugs); err == nil {
		// Create .malice folder in the users home directory
		er.CheckError(os.MkdirAll(maldirs.GetPluginsDir(), 0777))
		// Create the plugins config in the .malice folder
		er.CheckError(ioutil.WriteFile(configPath, tomlData, 0644))
		log.Debug("Malice plugins loaded from plugins/bindata.go")
	}
	er.CheckError(err)

	return
}

// func setInstalledFlag() {
// 	docker := client.NewDockerClient()
// 	for i, plugin := range Plugs.Plugins {
// 		if _, exists, _ := image.Exists(docker, plugin.Image); exists {
// 			Plugs.Plugins[i].Installed = true
// 		}
// 	}
// }
