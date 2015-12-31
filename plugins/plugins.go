package plugins

import (
	"fmt"

	"github.com/maliceio/malice/config"
)

// "github.com/pelletier/go-toml"

// ListEnabledPlugins lists all enabled plugins
func ListEnabledPlugins() {
	plugins := config.TConf.Plugins
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
