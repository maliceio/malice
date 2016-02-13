package plugins

import (
	"fmt"
	"strings"

	"github.com/crackcomm/go-clitable"
	"github.com/maliceio/malice/malice/maldocker"
)

// ListEnabledPlugins lists enabled plugins
func ListEnabledPlugins(detail bool) {
	// TODO: Create a template for this kind of output : http://stackoverflow.com/questions/10747054/special-case-treatment-for-the-last-element-of-a-range-in-google-gos-text-templ
	enabled := filterPluginsByEnabled()
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
	plugins := Plug.Plugins
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

	for _, plugin := range Plug.Plugins {
		if strings.EqualFold(plugin.Name, name) {
			return plugin
		}
	}

	return Plugin{}
}

// GetPluginsForMime will return all plugins that can consume the mime type file
func GetPluginsForMime(client *maldocker.Docker, mime string, installed bool) []Plugin {
	if installed {
		return filterPluginsByInstalled(client, mime)
	}
	return filterPluginsByMime(mime)
}

// filterPluginsByEnabled returns a map[string]plugin of enalbed plugins
func filterPluginsByInstalled(client *maldocker.Docker, mime string) []Plugin {
	installed := []Plugin{}

	for _, plugin := range filterPluginsByMime(mime) {
		if _, exists, _ := client.ImageExists(plugin.Image); exists {
			installed = append(installed, plugin)
		}
	}
	return installed
}

// filterPluginsByEnabled returns a map[string]plugin of plugins
// that work on the given mime type
func filterPluginsByMime(mime string) []Plugin {
	mimeMatch := []Plugin{}

	for _, plugin := range filterPluginsByEnabled() {
		if strings.Contains(plugin.Mime, mime) || strings.Contains(plugin.Mime, "*") {
			mimeMatch = append(mimeMatch, plugin)
		}
	}
	return mimeMatch
}

// filterPluginsByEnabled returns a map[string]plugin of enalbed plugins
func filterPluginsByEnabled() []Plugin {
	enabled := []Plugin{}

	for _, plugin := range Plug.Plugins {
		if plugin.Enabled {
			enabled = append(enabled, plugin)
		}
	}
	return enabled
}
