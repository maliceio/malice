package plugin

import "sync"

// Manager controls the plugin subsystem.
type Manager struct {
	config           ManagerConfig
	mu               sync.RWMutex // protects cMap
	muGC             sync.RWMutex // protects blobstore deletions
	cMap             map[*v2.Plugin]*controller
	containerdClient libcontainerd.Client
	blobStore        *basicBlobStore
	publisher        *pubsub.Publisher
}

// controller represents the manager's control on a plugin.
type controller struct {
	restart       bool
	exitChan      chan bool
	timeoutInSecs int
}
