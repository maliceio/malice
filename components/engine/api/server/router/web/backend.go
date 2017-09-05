package web

import (
	enginetypes "github.com/maliceio/engine/api/types"
	"github.com/maliceio/engine/api/types/filters"
)

// Backend for Plugin
type Backend interface {
	Start(name string, config *enginetypes.PluginDisableConfig) error
	Stop(name string, config *enginetypes.PluginEnableConfig) error
	BackUp(filters.Args) error
}
