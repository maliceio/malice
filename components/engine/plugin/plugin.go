package plugin

import (
	"sync"

	"github.com/maliceio/engine/api/types"
	digest "github.com/opencontainers/go-digest"
)

// Plugin represents an individual plugin.
type Plugin struct {
	mu        sync.RWMutex
	PluginObj types.Plugin `json:"plugin"` // todo: embed struct

	Config digest.Digest
}

// Name returns the plugin name.
func (p *Plugin) Name() string {
	return p.PluginObj.Name
}

// IsEnabled returns the active state of the plugin.
func (p *Plugin) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.PluginObj.Enabled
}

// GetID returns the plugin's ID.
func (p *Plugin) GetID() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.PluginObj.ID
}

// GetTypes returns the interface types of a plugin.
func (p *Plugin) GetTypes() []types.PluginInterfaceType {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.PluginObj.Config.Interface.Types
}
