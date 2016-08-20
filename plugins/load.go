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
func Load() {
	// Try to load plugins from
	// - git repo folder      : MALICE_ROOT/plugins/plugins.toml
	// - .malice folder       : $HOME/.malice/plugins.toml
	// - binary embedded file : bindata

	// Check for plugins config in repo
	if _, err := os.Stat("./plugins/plugins.toml"); err == nil {
		log.Debug("Malice plugins loaded from ./plugins/plugins.toml")

		_, err := toml.DecodeFile("./plugins/plugins.toml", &Plugs)
		er.CheckError(err)

		return
	}
	// Check for plugins config in .malice folder
	if _, err := os.Stat(path.Join(maldirs.GetBaseDir(), "./plugins.toml")); err == nil {
		homeConfigDir := path.Join(maldirs.GetBaseDir(), "./plugins.toml")
		log.Debug("Malice plugins loaded from ", homeConfigDir)

		_, err := toml.DecodeFile(homeConfigDir, &Plugs)
		er.CheckError(err)

		return
	}
	// Read plugin config out of bindata
	tomlData, err := Asset("plugins/plugins.toml")
	if err != nil {
		log.Error(err)
	}

	if _, err := toml.Decode(string(tomlData), &Plugs); err == nil {
		log.Debug("Malice plugins loaded from plugins/bindata.go")
		configPath := path.Join(maldirs.GetBaseDir(), "./plugins.toml")
		er.CheckError(ioutil.WriteFile(configPath, tomlData, 0644))
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
