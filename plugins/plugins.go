package plugins

import (
	"fmt"

	"github.com/maliceio/malice/config"
)

// "github.com/pelletier/go-toml"

// GetEnabledPlugins returns a map[string]plugin of enalbed plugins
func GetEnabledPlugins() map[string]config.Plugin {
	var enabled = make(map[string]config.Plugin)

	for name, plugin := range config.Conf.Plugins {
		if plugin.Enabled {
			enabled[name] = plugin
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

// ListAllPlugins lists all enabled plugins
func ListAllPlugins() {
	plugins := config.Conf.Plugins
	for name, plugin := range plugins {
		fmt.Println("Name: ", name)
		fmt.Println("Description: ", plugin.Description)
		fmt.Println("Image: ", plugin.Image)
		fmt.Println("Category: ", plugin.Category)
		fmt.Println("Mime: ", plugin.Mime)
		fmt.Println("---------------------")
	}
}
