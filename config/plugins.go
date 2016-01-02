package config

import "github.com/BurntSushi/toml"

// "github.com/pelletier/go-toml"

// PluginConfiguration represents the malice runtime plugins.
type PluginConfiguration struct {
	Plugins []Plugin `toml:"plugin"`
}

// Plugin represents a single plugin setting.
type Plugin struct {
	Name        string `toml:"name"`
	Enabled     bool   `toml:"enabled"`
	Category    string `toml:"category"`
	Description string `toml:"description"`
	Image       string `toml:"image"`
	Mime        string `toml:"mime"`
}

// Plug represents the Malice runtime configuration
var Plug PluginConfiguration

func init() {
	// Get the config file
	_, err := toml.DecodeFile("./plugins.toml", &Plug)
	assert(err)
	// fmt.Println(Plug)
}
