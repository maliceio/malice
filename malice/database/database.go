package database

import "github.com/maliceio/malice/plugins"

// GetPluginsByCategory gets malice plugins organized by category
func GetPluginsByCategory() map[string]interface{} {
	categoryList := make(map[string]interface{})
	for _, category := range plugins.GetCategories() {
		pluginList := make(map[string]interface{})
		for _, plugin := range plugins.GetAllPluginsInCategory(category) {
			pluginList[plugin.Name] = nil
		}
		categoryList[category] = pluginList
	}

	return categoryList
}

// GetPlugins gets malice plugins
func GetPlugins() map[string]interface{} {
	pluginList := make(map[string]interface{})
	for _, plugin := range plugins.Plugs.Plugins {
		pluginList[plugin.Name] = nil
	}

	return pluginList
}
