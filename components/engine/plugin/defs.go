package plugin

import (
	"sync"
)

// Store manages the plugin inventory in memory and on-disk
type Store struct {
	sync.RWMutex
	plugins map[string]*Plugin
}

// NewStore creates a Store.
func NewStore() *Store {
	return &Store{
		plugins: make(map[string]*Plugin),
	}
}

// CreateOpt is used to configure specific plugin details when created
type CreateOpt func(p *Plugin)
