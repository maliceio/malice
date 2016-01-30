package malice

import (
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/persist"
	"github.com/docker/machine/libmachine/ssh"
)

type API interface {
	persist.Store
	// persist.PluginDriverFactory
	NewHost(drivers.Driver) (*host.Host, error)
	Create(h *host.Host) error
}

type Client struct {
	*persist.PluginStore
	IsDebug        bool
	SSHClientType  ssh.ClientType
	GithubAPIToken string
}
