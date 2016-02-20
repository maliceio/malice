package plugins

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/data"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/maliceio/malice/malice/maldocker"
)

// Plugin represents a single plugin setting.
type Plugin struct {
	Name        string   `toml:"name"`
	Enabled     bool     `toml:"enabled"`
	Category    string   `toml:"category"`
	Description string   `toml:"description"`
	Image       string   `toml:"image"`
	APIKey      string   `toml:"apikey"`
	Mime        string   `toml:"mime"`
	HashTypes   []string `toml:"hashtypes"`
	Cmd         string   `toml:"cmd"`
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
	// Get the plugin file
	pluginPath := "./data/plugins.toml"
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		// er.CheckErrorNoStackWithMessage(err, "NOT FOUND")
		pluginPath = path.Join(maldirs.GetBaseDir(), "./plugins.toml")
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			pluginData, err := data.Asset("data/plugins.toml")
			er.CheckError(err)
			er.CheckError(ioutil.WriteFile(pluginPath, pluginData, 0644))
		}
	}
	log.Debug("Plugin Config: ", pluginPath)
	_, err := toml.DecodeFile(pluginPath, &Plugs)
	setInstalledFlag()
	// fmt.Printf("%#v\n", Plugs)
	er.CheckError(err)
}

func setInstalledFlag() {
	docker := maldocker.NewDockerClient()
	for i, plugin := range Plugs.Plugins {
		if _, exists, _ := docker.ImageExists(plugin.Image); exists {
			Plugs.Plugins[i].Installed = true
		}
	}
}
