package plugin

import "sync"

// ManagerConfig defines configuration needed to start new manager.
type ManagerConfig struct {
	// Store              *Store // remove
	// Executor           libcontainerd.Remote
	// RegistryService    registry.Service
	LiveRestoreEnabled bool // TODO: remove
	// LogPluginEvent     eventLogger
	Root     string
	ExecRoot string
	// AuthzMiddleware    *authorization.Middleware
}

// Manager controls the plugin subsystem.
type Manager struct {
	config ManagerConfig
	mu     sync.RWMutex // protects cMap
	muGC   sync.RWMutex // protects blobstore deletions
	// cMap             map[*v2.Plugin]*controller
	// containerdClient libcontainerd.Client
	// blobStore        *basicBlobStore
	// publisher        *pubsub.Publisher
}

// controller represents the manager's control on a plugin.
type controller struct {
	restart       bool
	exitChan      chan bool
	timeoutInSecs int
}
