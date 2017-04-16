package ui

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
)

// Start creates an Kibana container from the image blacktop/kibana:malice
func Start(docker *client.Docker, logs bool) (types.ContainerJSONBase, error) {

	portBindings := nat.PortMap{
		"5601/tcp": {{HostIP: "0.0.0.0", HostPort: "80"}},
	}

	if docker.Ping() {
		contJSON, err := container.Start(
			docker,                             // docker *client.Docker,
			nil,                                // cmd strslice.StrSlice,
			config.Conf.UI.Name,                // name string,
			config.Conf.UI.Image,               // image string,
			logs,                               // logs bool,
			nil,                                // binds []string,
			portBindings,                       // portBindings nat.PortMap,
			[]string{config.Conf.Docker.Links}, // links []string,
			nil, // env []string,
		)
		log.WithFields(log.Fields{
			"ip":   docker.GetIP(),
			"port": config.Conf.UI.Ports,
			"name": contJSON.Name,
			"env":  config.Conf.Environment.Run,
		}).Info("Kibana Container Started")

		return contJSON, err
	}
	return types.ContainerJSONBase{}, errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
}

// Init initalizes Kibana for use with malice
func Init(addr string) error {

	var err error

	return err
}
