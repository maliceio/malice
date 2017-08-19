package plugin

import (
	"sync"

	"github.com/maliceio/engine/api/types"
)

// Plugin represents an individual plugin.
type Plugin struct {
	mu        sync.RWMutex
	PluginObj types.Plugin `json:"plugin"` // todo: embed struct

	Rootfs string

	Config string

	SwarmServiceID string
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
