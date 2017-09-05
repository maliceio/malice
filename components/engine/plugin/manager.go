package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/pkg/pubsub"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/registry"
	"github.com/docker/swarmkit/ioutils"
	"github.com/maliceio/engine/api/types"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

const configFileName = "config.json"

var validFullID = regexp.MustCompile(`^([a-f0-9]{64})$`)

// ManagerConfig defines configuration needed to start new manager.
type ManagerConfig struct {
	Store *Store // remove
	// Executor           libcontainerd.Remote
	RegistryService    registry.Service
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
	cMap   map[*Plugin]*controller
	// containerdClient libcontainerd.Client
	// blobStore        *basicBlobStore
	publisher *pubsub.Publisher
}

// controller represents the manager's control on a plugin.
type controller struct {
	restart       bool
	exitChan      chan bool
	timeoutInSecs int
}

// pluginRegistryService ensures that all resolved repositories
// are of the plugin class.
type pluginRegistryService struct {
	registry.Service
}

func (s pluginRegistryService) ResolveRepository(name reference.Named) (repoInfo *registry.RepositoryInfo, err error) {
	repoInfo, err = s.Service.ResolveRepository(name)
	if repoInfo != nil {
		repoInfo.Class = "plugin"
	}
	return
}

// NewManager returns a new plugin manager.
func NewManager(config ManagerConfig) (*Manager, error) {
	if config.RegistryService != nil {
		config.RegistryService = pluginRegistryService{config.RegistryService}
	}
	manager := &Manager{
		config: config,
	}

	manager.cMap = make(map[*Plugin]*controller)
	if err := manager.reload(); err != nil {
		return nil, errors.Wrap(err, "failed to restore plugins")
	}

	manager.publisher = pubsub.NewPublisher(0, 0)
	return manager, nil
}

func handleLoadError(err error, id string) {
	if err == nil {
		return
	}
	logger := logrus.WithError(err).WithField("id", id)
	if os.IsNotExist(errors.Cause(err)) {
		// Likely some error while removing on an older version of docker
		logger.Warn("missing plugin config, skipping: this may be caused due to a failed remove and requires manual cleanup.")
		return
	}
	logger.Error("error loading plugin, skipping")
}

func (pm *Manager) reload() error {
	dir, err := ioutil.ReadDir(pm.config.Root)
	if err != nil {
		return errors.Wrapf(err, "failed to read %v", pm.config.Root)
	}
	plugins := make(map[string]*Plugin)
	for _, v := range dir {
		if validFullID.MatchString(v.Name()) {
			p, err := pm.loadPlugin(v.Name())
			if err != nil {
				handleLoadError(err, v.Name())
				continue
			}
			plugins[p.GetID()] = p
		} else {
			if validFullID.MatchString(strings.TrimSuffix(v.Name(), "-removing")) {
				// There was likely some error while removing this plugin, let's try to remove again here
				if err := system.EnsureRemoveAll(v.Name()); err != nil {
					logrus.WithError(err).WithField("id", v.Name()).Warn("error while attempting to clean up previously removed plugin")
				}
			}
		}
	}

	pm.config.Store.SetAll(plugins)

	var wg sync.WaitGroup
	wg.Add(len(plugins))
	for _, p := range plugins {
		c := &controller{} // todo: remove this
		pm.cMap[p] = c
		go func(p *Plugin) {
			defer wg.Done()

			pm.save(p)
			requiresManualRestore := !pm.config.LiveRestoreEnabled && p.IsEnabled()

			if requiresManualRestore {
				// if liveRestore is not enabled, the plugin will be stopped now so we should enable it
				if err := pm.enable(p, c, true); err != nil {
					logrus.Errorf("failed to enable plugin '%s': %s", p.Name(), err)
				}
			}
		}(p)
	}
	wg.Wait()
	return nil
}

// Get looks up the requested plugin in the store.
func (pm *Manager) Get(idOrName string) (*Plugin, error) {
	return pm.config.Store.GetPlugin(idOrName)
}

func (pm *Manager) loadPlugin(id string) (*Plugin, error) {
	p := filepath.Join(pm.config.Root, id, configFileName)
	dt, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading %v", p)
	}
	var plugin Plugin
	if err := json.Unmarshal(dt, &plugin); err != nil {
		return nil, errors.Wrapf(err, "error decoding %v", p)
	}
	return &plugin, nil
}

func (pm *Manager) save(p *Plugin) error {
	pluginJSON, err := json.Marshal(p)
	if err != nil {
		return errors.Wrap(err, "failed to marshal plugin json")
	}
	if err := ioutils.AtomicWriteFile(filepath.Join(pm.config.Root, p.GetID(), configFileName), pluginJSON, 0600); err != nil {
		return errors.Wrap(err, "failed to write atomically plugin json")
	}
	return nil
}

func (pm *Manager) enable(p *Plugin, c *controller, force bool) error {
	return nil
}

func shutdownPlugin(p *Plugin, c *controller) {
	// pluginID := p.GetID()

	// err := containerdClient.Signal(pluginID, int(unix.SIGTERM))
	// if err != nil {
	// 	logrus.Errorf("Sending SIGTERM to plugin failed with error: %v", err)
	// } else {
	select {
	case <-c.exitChan:
		logrus.Debug("Clean shutdown of plugin")
	case <-time.After(time.Second * 10):
		logrus.Debug("Force shutdown plugin")
		// if err := containerdClient.Signal(pluginID, int(unix.SIGKILL)); err != nil {
		// 	logrus.Errorf("Sending SIGKILL to plugin failed with error: %v", err)
		// }
	}
	// }
}

func (pm *Manager) disable(p *Plugin, c *controller) error {
	if !p.IsEnabled() {
		return errors.Wrap(errDisabled(p.Name()), "plugin is already disabled")
	}

	c.restart = false
	shutdownPlugin(p, c)
	pm.config.Store.SetState(p, false)
	return pm.save(p)
}

// Shutdown stops all plugins and called during daemon shutdown.
func (pm *Manager) Shutdown() {
	plugins := pm.config.Store.GetAll()
	for _, p := range plugins {
		pm.mu.RLock()
		c := pm.cMap[p]
		pm.mu.RUnlock()

		if p.IsEnabled() {
			c.restart = false
			shutdownPlugin(p, c)
		}
	}
}

func (pm *Manager) setupNewPlugin(configDigest digest.Digest) (types.PluginConfig, error) {
	var config types.PluginConfig
	return config, nil
}

func (pm *Manager) upgradePlugin(p *Plugin, configDigest digest.Digest) (err error) {
	// config, err := pm.setupNewPlugin(configDigest)
	// if err != nil {
	// 	return err
	// }

	defer func() {
		// cleanup
		fmt.Println("perform any cleanup after upgrade")
	}()

	// TODO: add docker code to pull updated version of the image

	return nil
}
