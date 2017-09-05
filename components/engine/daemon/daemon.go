package daemon

import (
	"github.com/Sirupsen/logrus"
	"github.com/maliceio/engine/daemon/config"
	"github.com/maliceio/engine/malice/version"
	"github.com/maliceio/engine/plugin"
)

// Daemon holds information about the Docker daemon.
type Daemon struct {
	ID         string
	repository string

	configStore *config.Config

	root            string
	seccompEnabled  bool
	apparmorEnabled bool
	shutdown        bool

	PluginStore   *plugin.Store // todo: remove
	pluginManager *plugin.Manager

	machineMemory uint64

	seccompProfile     []byte
	seccompProfilePath string

	hosts       map[string]bool // hosts stores the addresses the daemon is listening on
	startupDone chan struct{}
}

// NewDaemon sets up everything for the daemon to be able to service
// requests from the webserver.
func NewDaemon(config *config.Config) (daemon *Daemon, err error) {

	// Validate platform-specific requirements
	if err := checkSystem(); err != nil {
		return nil, err
	}

	d := &Daemon{
		configStore: config,
		startupDone: make(chan struct{}),
	}
	// Ensure the daemon is properly shutdown if there is a failure during
	// initialization
	defer func() {
		if err != nil {
			if err := d.Shutdown(); err != nil {
				logrus.Error(err)
			}
		}
	}()

	// // Plugin system initialization should happen before restore. Do not change order.
	// d.pluginManager, err = plugin.NewManager(plugin.ManagerConfig{
	// 	Root:               filepath.Join(config.Root, "plugins"),
	// 	ExecRoot:           getPluginExecRoot(config.Root),
	// 	Store:              d.PluginStore,
	// 	Executor:           containerdRemote,
	// 	RegistryService:    registryService,
	// 	LiveRestoreEnabled: config.LiveRestoreEnabled,
	// 	LogPluginEvent:     d.LogPluginEvent, // todo: make private
	// 	AuthzMiddleware:    config.AuthzMiddleware,
	// })
	// if err != nil {
	// 	return nil, errors.Wrap(err, "couldn't create plugin manager")
	// }

	// // Configure the volumes driver
	// volStore, err := d.configureVolumes(rootIDs)
	// if err != nil {
	// 	return nil, err
	// }

	// trustKey, err := api.LoadOrCreateTrustKey(config.TrustKeyPath)
	// if err != nil {
	// 	return nil, err
	// }

	logrus.WithFields(logrus.Fields{
		"version": version.Version,
		"commit":  version.GitCommit,
	}).Info("Malice daemon")

	return d, nil
}

func (daemon *Daemon) waitForStartupDone() {
	<-daemon.startupDone
}

// Shutdown stops the daemon.
func (daemon *Daemon) Shutdown() error {
	daemon.shutdown = true
}

// IsShuttingDown tells whether the daemon is shutting down or not
func (daemon *Daemon) IsShuttingDown() bool {
	return daemon.shutdown
}

// func (daemon *Daemon) configureVolumes(rootIDs idtools.IDPair) (*store.VolumeStore, error) {
// 	volumesDriver, err := local.New(daemon.configStore.Root, rootIDs)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	volumedrivers.RegisterPluginGetter(daemon.PluginStore)
//
// 	if !volumedrivers.Register(volumesDriver, volumesDriver.Name()) {
// 		return nil, errors.New("local volume driver could not be registered")
// 	}
// 	return store.New(daemon.configStore.Root)
// }
