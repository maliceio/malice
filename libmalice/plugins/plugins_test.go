package plugins

import (
	"fmt"
	"testing"

	"github.com/maliceio/malice/config"
)

// "github.com/pelletier/go-toml"

// GetEnabledPlugins returns a map[string]plugin of enalbed plugins
func TestGetEnabledPlugins(t *testing.T) {
	var enabled = GetEnabledPlugins()
	fmt.Println(enabled)
}

func TestListEnabledPlugins(t *testing.T) {
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
