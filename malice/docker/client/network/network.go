package network

import (
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/network"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	"golang.org/x/net/context"

	"regexp"

	log "github.com/Sirupsen/logrus"
)

// Exists returns type.NetworkResource and true
// if the network name exists, otherwise false.
func Exists(docker *client.Docker, name string) (types.NetworkResource, bool, error) {
	return parseNetworks(docker, name, true)
}

// Create creates a docker Network with the given name
// returns: NetworkCreateResponse, error
func Create(docker *client.Docker, name string) (types.NetworkCreateResponse, error) {
	options := types.NetworkCreate{
	// true,           // CheckDuplicate bool
	// "bridge",       // Driver         string
	// false,          // EnableIPv6     bool
	// network.IPAM{}, // IPAM           network.IPAM
	// false,          // Internal       bool
	// nil,            // Options        map[string]string
	// nil,            // Labels         map[string]string
	}
	net, err := docker.Client.NetworkCreate(context.Background(), name, options)
	log.WithFields(log.Fields{
		"name": name,
		"env":  config.Conf.Environment.Run,
	}).Info("Created Network: ", name)
	return net, err
}

// Connect connects a container to a network
func Connect(docker *client.Docker, net types.NetworkResource, container types.ContainerJSONBase) error {
	netConfig := network.EndpointSettings{}
	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Debugf("Connecting container %s to network %s", container.Name, net.Name)
	return docker.Client.NetworkConnect(context.Background(), net.ID, container.ID, &netConfig)
}

// parseNetworks parses the networks
func parseNetworks(docker *client.Docker, name string, all bool) (types.NetworkResource, bool, error) {
	// list networks
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Searching for Network: ", name)
	networks, err := List(docker, all)
	if err != nil {
		return types.NetworkResource{}, false, err
	}
	// locate docker Network that matches name
	r := regexp.MustCompile(name)
	if len(networks) != 0 {
		for _, network := range networks {
			if r.MatchString(network.Name) {
				log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Network FOUND: ", name)
				return network, true, nil
			}
		}
	}
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Network NOT Found: ", name)
	return types.NetworkResource{}, false, nil
}

// List returns array of type NetworkResources and error
func List(docker *client.Docker, all bool) ([]types.NetworkResource, error) {

	options := types.NetworkListOptions{Filters: filters.Args{}}
	networks, err := docker.Client.NetworkList(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return networks, nil
}
