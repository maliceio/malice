package plugin

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/system"
	"github.com/maliceio/engine/api/types"
	"github.com/pkg/errors"
)

// Disable deactivates a plugin. This means resources (volumes, networks) cant use them.
func (pm *Manager) Disable(refOrID string, config *types.PluginDisableConfig) error {
	p, err := pm.config.Store.GetV2Plugin(refOrID)
	if err != nil {
		return err
	}
	pm.mu.RLock()
	c := pm.cMap[p]
	pm.mu.RUnlock()

	if !config.ForceDisable && p.GetRefCount() > 0 {
		return errors.WithStack(inUseError(p.Name()))
	}

	for _, typ := range p.GetTypes() {
		if typ.Capability == authorization.AuthZApiImplements {
			pm.config.AuthzMiddleware.RemovePlugin(p.Name())
		}
	}

	if err := pm.disable(p, c); err != nil {
		return err
	}
	pm.publisher.Publish(EventDisable{Plugin: p.PluginObj})
	pm.config.LogPluginEvent(p.GetID(), refOrID, "disable")
	return nil
}

// Enable activates a plugin, which implies that they are ready to be used by containers.
func (pm *Manager) Enable(refOrID string, config *types.PluginEnableConfig) error {
	p, err := pm.config.Store.GetV2Plugin(refOrID)
	if err != nil {
		return err
	}

	c := &controller{timeoutInSecs: config.Timeout}
	if err := pm.enable(p, c, false); err != nil {
		return err
	}
	pm.publisher.Publish(EventEnable{Plugin: p.PluginObj})
	pm.config.LogPluginEvent(p.GetID(), refOrID, "enable")
	return nil
}

// Upgrade upgrades a plugin
func (pm *Manager) Upgrade(ctx context.Context, ref reference.Named, name string, metaHeader http.Header, authConfig *types.AuthConfig, privileges types.PluginPrivileges, outStream io.Writer) (err error) {
	p, err := pm.config.Store.GetV2Plugin(name)
	if err != nil {
		return err
	}

	if p.IsEnabled() {
		return errors.Wrap(enabledError(p.Name()), "plugin must be disabled before upgrading")
	}

	pm.muGC.RLock()
	defer pm.muGC.RUnlock()

	// revalidate because Pull is public
	if _, err := reference.ParseNormalizedNamed(name); err != nil {
		return errors.Wrapf(validationError{err}, "failed to parse %q", name)
	}

	tmpRootFSDir, err := ioutil.TempDir(pm.tmpDir(), ".rootfs")
	if err != nil {
		return errors.Wrap(systemError{err}, "error preparing upgrade")
	}
	defer os.RemoveAll(tmpRootFSDir)

	dm := &downloadManager{
		tmpDir:    tmpRootFSDir,
		blobStore: pm.blobStore,
	}

	pluginPullConfig := &distribution.ImagePullConfig{
		Config: distribution.Config{
			MetaHeaders:      metaHeader,
			AuthConfig:       authConfig,
			RegistryService:  pm.config.RegistryService,
			ImageEventLogger: pm.config.LogPluginEvent,
			ImageStore:       dm,
		},
		DownloadManager: dm, // todo: reevaluate if possible to substitute distribution/xfer dependencies instead
		Schema2Types:    distribution.PluginTypes,
	}

	err = pm.pull(ctx, ref, pluginPullConfig, outStream)
	if err != nil {
		go pm.GC()
		return err
	}

	if err := pm.upgradePlugin(p, dm.configDigest, dm.blobs, tmpRootFSDir, &privileges); err != nil {
		return err
	}
	p.PluginObj.PluginReference = ref.String()
	return nil
}

// Pull pulls a plugin, check if the correct privileges are provided and install the plugin.
func (pm *Manager) Pull(ctx context.Context, ref reference.Named, name string, metaHeader http.Header, authConfig *types.AuthConfig, privileges types.PluginPrivileges, outStream io.Writer, opts ...CreateOpt) (err error) {
	pm.muGC.RLock()
	defer pm.muGC.RUnlock()

	// revalidate because Pull is public
	nameref, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return errors.Wrapf(validationError{err}, "failed to parse %q", name)
	}
	name = reference.FamiliarString(reference.TagNameOnly(nameref))

	if err := pm.config.Store.validateName(name); err != nil {
		return validationError{err}
	}

	tmpRootFSDir, err := ioutil.TempDir(pm.tmpDir(), ".rootfs")
	if err != nil {
		return errors.Wrap(systemError{err}, "error preparing pull")
	}
	defer os.RemoveAll(tmpRootFSDir)

	dm := &downloadManager{
		tmpDir:    tmpRootFSDir,
		blobStore: pm.blobStore,
	}

	pluginPullConfig := &distribution.ImagePullConfig{
		Config: distribution.Config{
			MetaHeaders:      metaHeader,
			AuthConfig:       authConfig,
			RegistryService:  pm.config.RegistryService,
			ImageEventLogger: pm.config.LogPluginEvent,
			ImageStore:       dm,
		},
		DownloadManager: dm, // todo: reevaluate if possible to substitute distribution/xfer dependencies instead
		Schema2Types:    distribution.PluginTypes,
	}

	err = pm.pull(ctx, ref, pluginPullConfig, outStream)
	if err != nil {
		go pm.GC()
		return err
	}

	refOpt := func(p *v2.Plugin) {
		p.PluginObj.PluginReference = ref.String()
	}
	optsList := make([]CreateOpt, 0, len(opts)+1)
	optsList = append(optsList, opts...)
	optsList = append(optsList, refOpt)

	p, err := pm.createPlugin(name, dm.configDigest, dm.blobs, tmpRootFSDir, &privileges, optsList...)
	if err != nil {
		return err
	}

	pm.publisher.Publish(EventCreate{Plugin: p.PluginObj})
	return nil
}

// List displays the list of plugins and associated metadata.
func (pm *Manager) List(pluginFilters filters.Args) ([]types.Plugin, error) {
	if err := pluginFilters.Validate(acceptedPluginFilterTags); err != nil {
		return nil, err
	}

	enabledOnly := false
	disabledOnly := false
	if pluginFilters.Include("enabled") {
		if pluginFilters.ExactMatch("enabled", "true") {
			enabledOnly = true
		} else if pluginFilters.ExactMatch("enabled", "false") {
			disabledOnly = true
		} else {
			return nil, invalidFilter{"enabled", pluginFilters.Get("enabled")}
		}
	}

	plugins := pm.config.Store.GetAll()
	out := make([]types.Plugin, 0, len(plugins))

next:
	for _, p := range plugins {
		if enabledOnly && !p.PluginObj.Enabled {
			continue
		}
		if disabledOnly && p.PluginObj.Enabled {
			continue
		}
		if pluginFilters.Include("capability") {
			for _, f := range p.GetTypes() {
				if !pluginFilters.Match("capability", f.Capability) {
					continue next
				}
			}
		}
		out = append(out, p.PluginObj)
	}
	return out, nil
}

// Push pushes a plugin to the store.
func (pm *Manager) Push(ctx context.Context, name string, metaHeader http.Header, authConfig *types.AuthConfig, outStream io.Writer) error {
	p, err := pm.config.Store.GetV2Plugin(name)
	if err != nil {
		return err
	}

	ref, err := reference.ParseNormalizedNamed(p.Name())
	if err != nil {
		return errors.Wrapf(err, "plugin has invalid name %v for push", p.Name())
	}

	var po progress.Output
	if outStream != nil {
		// Include a buffer so that slow client connections don't affect
		// transfer performance.
		progressChan := make(chan progress.Progress, 100)

		writesDone := make(chan struct{})

		defer func() {
			close(progressChan)
			<-writesDone
		}()

		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithCancel(ctx)

		go func() {
			progressutils.WriteDistributionProgress(cancelFunc, outStream, progressChan)
			close(writesDone)
		}()

		po = progress.ChanOutput(progressChan)
	} else {
		po = progress.DiscardOutput()
	}

	// TODO: replace these with manager
	is := &pluginConfigStore{
		pm:     pm,
		plugin: p,
	}
	ls := &pluginLayerProvider{
		pm:     pm,
		plugin: p,
	}
	rs := &pluginReference{
		name:     ref,
		pluginID: p.Config,
	}

	uploadManager := xfer.NewLayerUploadManager(3)

	imagePushConfig := &distribution.ImagePushConfig{
		Config: distribution.Config{
			MetaHeaders:      metaHeader,
			AuthConfig:       authConfig,
			ProgressOutput:   po,
			RegistryService:  pm.config.RegistryService,
			ReferenceStore:   rs,
			ImageEventLogger: pm.config.LogPluginEvent,
			ImageStore:       is,
			RequireSchema2:   true,
		},
		ConfigMediaType: schema2.MediaTypePluginConfig,
		LayerStore:      ls,
		UploadManager:   uploadManager,
	}

	return distribution.Push(ctx, ref, imagePushConfig)
}

// Remove deletes plugin's root directory.
func (pm *Manager) Remove(name string, config *types.PluginRmConfig) error {
	p, err := pm.config.Store.GetV2Plugin(name)
	pm.mu.RLock()
	c := pm.cMap[p]
	pm.mu.RUnlock()

	if err != nil {
		return err
	}

	if !config.ForceRemove {
		if p.GetRefCount() > 0 {
			return inUseError(p.Name())
		}
		if p.IsEnabled() {
			return enabledError(p.Name())
		}
	}

	if p.IsEnabled() {
		if err := pm.disable(p, c); err != nil {
			logrus.Errorf("failed to disable plugin '%s': %s", p.Name(), err)
		}
	}

	defer func() {
		go pm.GC()
	}()

	id := p.GetID()
	pluginDir := filepath.Join(pm.config.Root, id)

	if err := mount.RecursiveUnmount(pluginDir); err != nil {
		return errors.Wrap(err, "error unmounting plugin data")
	}

	removeDir := pluginDir + "-removing"
	if err := os.Rename(pluginDir, removeDir); err != nil {
		return errors.Wrap(err, "error performing atomic remove of plugin dir")
	}

	if err := system.EnsureRemoveAll(removeDir); err != nil {
		return errors.Wrap(err, "error removing plugin dir")
	}
	pm.config.Store.Remove(p)
	pm.config.LogPluginEvent(id, name, "remove")
	pm.publisher.Publish(EventRemove{Plugin: p.PluginObj})
	return nil
}
