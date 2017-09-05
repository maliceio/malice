package plugin

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
)

// GetPlugin retrieves a plugin by name, id or partial ID.
func (ps *Store) GetPlugin(refOrID string) (*Plugin, error) {
	ps.RLock()
	defer ps.RUnlock()

	id, err := ps.resolvePluginID(refOrID)
	if err != nil {
		return nil, err
	}

	p, idOk := ps.plugins[id]
	if !idOk {
		return nil, errors.WithStack(errNotFound(id))
	}

	return p, nil
}

// validateName returns error if name is already reserved. always call with lock and full name
func (ps *Store) validateName(name string) error {
	for _, p := range ps.plugins {
		if p.Name() == name {
			return alreadyExistsError(name)
		}
	}
	return nil
}

// GetAll retrieves all plugins.
func (ps *Store) GetAll() map[string]*Plugin {
	ps.RLock()
	defer ps.RUnlock()
	return ps.plugins
}

// SetAll initialized plugins during daemon restore.
func (ps *Store) SetAll(plugins map[string]*Plugin) {
	ps.Lock()
	defer ps.Unlock()
	ps.plugins = plugins
}

// SetState sets the active state of the plugin and updates plugindb.
func (ps *Store) SetState(p *Plugin, state bool) {
	ps.Lock()
	defer ps.Unlock()

	p.PluginObj.Enabled = state
}

// Add adds a plugin to memory and plugindb.
// An error will be returned if there is a collision.
func (ps *Store) Add(p *Plugin) error {
	ps.Lock()
	defer ps.Unlock()

	if v, exist := ps.plugins[p.GetID()]; exist {
		return fmt.Errorf("plugin %q has the same ID %s as %q", p.Name(), p.GetID(), v.Name())
	}
	ps.plugins[p.GetID()] = p
	return nil
}

// Remove removes a plugin from memory and plugindb.
func (ps *Store) Remove(p *Plugin) {
	ps.Lock()
	delete(ps.plugins, p.GetID())
	ps.Unlock()
}

func (ps *Store) resolvePluginID(idOrName string) (string, error) {
	ps.RLock() // todo: fix
	defer ps.RUnlock()

	if validFullID.MatchString(idOrName) {
		return idOrName, nil
	}

	ref, err := reference.ParseNormalizedNamed(idOrName)
	if err != nil {
		return "", errors.WithStack(errNotFound(idOrName))
	}
	if _, ok := ref.(reference.Canonical); ok {
		logrus.Warnf("canonical references cannot be resolved: %v", reference.FamiliarString(ref))
		return "", errors.WithStack(errNotFound(idOrName))
	}

	ref = reference.TagNameOnly(ref)

	for _, p := range ps.plugins {
		if p.PluginObj.Name == reference.FamiliarString(ref) {
			return p.PluginObj.ID, nil
		}
	}

	var found *Plugin
	for id, p := range ps.plugins { // this can be optimized
		if strings.HasPrefix(id, idOrName) {
			if found != nil {
				return "", errors.WithStack(errAmbiguous(idOrName))
			}
			found = p
		}
	}
	if found == nil {
		return "", errors.WithStack(errNotFound(idOrName))
	}
	return found.PluginObj.ID, nil
}
