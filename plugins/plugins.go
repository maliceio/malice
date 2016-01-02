package plugins

import (
	"fmt"

	"github.com/maliceio/malice/config"
)

// "github.com/pelletier/go-toml"

// GetEnabledPlugins returns a map[string]plugin of enalbed plugins
func GetEnabledPlugins() map[string]tomlConfig.Plugins {
	var enabled map[string]tomlConfig.Plugins
	plugins := config.TConf.Plugins
	for name, plug := range plugins {
		if plug.Enabled {
			enabled[name] = append(enabled[name], plug)
		}
	}
	return enabled
}

// ListEnabledPlugins lists all enabled plugins
func ListEnabledPlugins() {
	plugins := config.Conf.Plugins
	for name, plugin := range plugins {
		if plugin.Enabled {
			fmt.Println("Name: ", name)
			fmt.Println("Description: ", plugin.Description)
			fmt.Println("Image: ", plugin.Image)
			fmt.Println("Category: ", plugin.Category)
			fmt.Println("Mime: ", plugin.Mime)
		}
	}
}
