package plugins

import (
	"fmt"
	"strings"

	"github.com/crackcomm/go-clitable"
	"github.com/maliceio/malice/utils"
)

// ListEnabledPlugins lists enabled plugins
func ListEnabledPlugins(detail bool) {
	// TODO: Create a template for this kind of output : http://stackoverflow.com/questions/10747054/special-case-treatment-for-the-last-element-of-a-range-in-google-gos-text-templ
	enabled := getEnabled(Plugs.Plugins)
	if detail {
		ToMarkDownTable(enabled)
	} else {
		for _, plugin := range enabled {
			fmt.Println(plugin.Name)
		}
	}
}

// ListAllPlugins lists all plugins
func ListAllPlugins(detail bool) {
	plugins := Plugs.Plugins
	if detail {
		ToMarkDownTable(plugins)
	} else {
		for _, plugin := range plugins {
			fmt.Println(plugin.Name)
		}
	}
}

// ToMarkDownTable prints plugins out as Markdown table
func ToMarkDownTable(plugins []Plugin) {
	table := clitable.New([]string{"Name", "Description", "Enabled", "Image", "Category", "Mime"})
	for _, plugin := range plugins {
		table.AddRow(map[string]interface{}{
			"Name":        plugin.Name,
			"Description": plugin.Description,
			"Enabled":     plugin.Enabled,
			"Image":       plugin.Image,
			"Category":    plugin.Category,
			"Mime":        plugin.Mime,
		})
	}
	table.Markdown = true
	table.Print()
}

// GetPluginByName will return plugin for the given name
func GetPluginByName(name string) Plugin {
	for _, plugin := range Plugs.Plugins {
		if strings.EqualFold(plugin.Name, name) {
			return plugin
		}
	}
	return Plugin{}
}

// GetIntelPlugins will return all Intel plugins
func GetIntelPlugins(hashType string, enabled bool) []Plugin {
	var intelPlugs []Plugin
	if enabled {
		intelPlugs = getIntel(getEnabled(getInstalled()))
	} else {
		intelPlugs = getIntel(getInstalled())
	}
	// filter down to intel plugins with apikey's set in ENV
	// var allSet bool
	// var hasEnvPlugs []Plugin
	// for _, plugin := range intelPlugs {
	// 	allSet = true
	// 	for _, pluginEnv := range plugin.Env {
	// 		if os.Getenv(pluginEnv) == "" {
	// 			allSet = false
	// 		}
	// 	}
	// 	if allSet {
	// 		if utils.StringInSlice(hashType, plugin.HashTypes) {
	// 			hasEnvPlugs = append(hasEnvPlugs, plugin)
	// 		}
	// 	}
	// }
	return intelPlugs
}

// GetPluginsForMime will return all plugins that can consume the mime type file
func GetPluginsForMime(mime string, enabled bool) []Plugin {
	if enabled {
		return getMime(mime, getEnabled(getInstalled()))
	}
	return getMime(mime, getInstalled())
}

func getIntel(plugins []Plugin) []Plugin {
	intel := []Plugin{}
	if plugins == nil {
		plugins = Plugs.Plugins
	}
	for _, plugin := range plugins {
		if strings.Contains(plugin.Category, "intel") {
			intel = append(intel, plugin)
		}
	}
	return intel
}

// getInstalled returns a map[string]plugin of installed plugins
func getInstalled() []Plugin {
	installed := []Plugin{}
	for _, plugin := range Plugs.Plugins {
		if plugin.Installed {
			installed = append(installed, plugin)
		}
	}
	return installed
}

// GetCategories returns all categories
func GetCategories() []string {
	categories := []string{}
	for _, plugin := range Plugs.Plugins {
		if !utils.StringInSlice(plugin.Category, categories) {
			categories = append(categories, plugin.Category)
		}
	}
	return categories
}

// GetAllPluginsInCategory returns all plugins in a give category
func GetAllPluginsInCategory(category string) []Plugin {
	inCategory := []Plugin{}
	for _, plugin := range Plugs.Plugins {
		if strings.EqualFold(plugin.Category, category) {
			inCategory = append(inCategory, plugin)
		}
	}
	return inCategory
}

// GetEnabledPlugins will return all enabled plugins
func GetEnabledPlugins() []Plugin {
	return getEnabled(Plugs.Plugins)
}

// filterPluginsByEnabled returns a map[string]plugin of plugins
// that work on the given mime type
func getMime(mime string, plugins []Plugin) []Plugin {
	mimeMatch := []Plugin{}
	if plugins == nil {
		plugins = Plugs.Plugins
	}
	for _, plugin := range plugins {
		if strings.Contains(plugin.Mime, mime) || strings.Contains(plugin.Mime, "*") {
			mimeMatch = append(mimeMatch, plugin)
		}
	}
	return mimeMatch
}

// getEnabled returns a map[string]plugin of enabled plugins
func getEnabled(plugins []Plugin) []Plugin {
	enabled := []Plugin{}
	if plugins == nil {
		return Plugs.Plugins
	}
	for _, plugin := range Plugs.Plugins {
		if plugin.Enabled {
			enabled = append(enabled, plugin)
		}
	}
	return enabled
}
