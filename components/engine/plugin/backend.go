package plugin

import (
	"context"

	"github.com/maliceio/engine/api/types"
)

// Disable deactivates a plugin. This means resources (volumes, networks) cant use them.
func (pm *Manager) Disable(refOrID string, config *types.PluginDisableConfig) error {
	return nil
}

// Enable activates a plugin, which implies that they are ready to be used by containers.
func (pm *Manager) Enable(refOrID string, config *types.PluginEnableConfig) error {
	return nil
}

// Upgrade upgrades a plugin
func (pm *Manager) Upgrade(ctx context.Context, name string) (err error) {
	return nil
}

// Pull pulls a plugin, check if the correct privileges are provided and install the plugin.
func (pm *Manager) Pull(ctx context.Context, name string) (err error) {
	pm.muGC.RLock()
	defer pm.muGC.RUnlock()

	return nil
}

// List displays the list of plugins and associated metadata.
func (pm *Manager) List() ([]types.Plugin, error) {
	// if err := pluginFilters.Validate(acceptedPluginFilterTags); err != nil {
	// 	return nil, err
	// }

	// enabledOnly := false
	// disabledOnly := false
	// if pluginFilters.Include("enabled") {
	// 	if pluginFilters.ExactMatch("enabled", "true") {
	// 		enabledOnly = true
	// 	} else if pluginFilters.ExactMatch("enabled", "false") {
	// 		disabledOnly = true
	// 	} else {
	// 		return nil, invalidFilter{"enabled", pluginFilters.Get("enabled")}
	// 	}
	// }

	plugins := pm.config.Store.GetAll()
	out := make([]types.Plugin, 0, len(plugins))

	// for _, p := range plugins {
	// 	if enabledOnly && !p.PluginObj.Enabled {
	// 		continue
	// 	}
	// 	if disabledOnly && p.PluginObj.Enabled {
	// 		continue
	// 	}
	// 	out = append(out, p.PluginObj)
	// }
	return out, nil
}

// Push pushes a plugin to the store.
func (pm *Manager) Push(ctx context.Context, name string) error {
	return nil
}

// Remove deletes plugin's root directory.
func (pm *Manager) Remove(name string, config *types.PluginRmConfig) error {
	return nil
}
