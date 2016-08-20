package container

import (
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	er "github.com/maliceio/malice/malice/errors"
)

// List returns array of types.Containers and error
func List(docker *client.Docker, all bool) ([]types.Container, error) {
	options := types.ContainerListOptions{
		All: true,
		// Limit:  opts.last,
		// Size:   opts.size,
		// Filter: containerFilterArgs,
	}
	containers, err := docker.Client.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// Inspect returns types.ContainerJSON from Container ID
// if the container name exists, otherwise false.
func Inspect(docker *client.Docker, id string) (types.ContainerJSON, error) {
	contJSON, err := docker.Client.ContainerInspect(context.Background(), id)
	return contJSON, err
}

// Exists returns APIContainers containers list and true
// if the container name exists, otherwise false.
func Exists(docker *client.Docker, name string) (types.Container, bool, error) {
	return parseContainers(docker, strings.TrimLeft(name, "/"), true)
}

// Running returns APIContainers containers list and true
// if the container name exists and is running, otherwise false.
func Running(docker *client.Docker, name string) (types.Container, bool, error) {
	return parseContainers(docker, strings.TrimLeft(name, "/"), false)
}

func parseContainers(docker *client.Docker, name string, all bool) (types.Container, bool, error) {
	// list containers
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Searching for container: ", name)
	containers, err := List(docker, all)
	if err != nil {
		return types.Container{}, false, err
	}
	// locate docker container that matches name
	if len(containers) != 0 {
		for _, container := range containers {

			cont, err := Inspect(docker, container.ID)
			er.CheckError(err)

			log.Debugln("name: ", name, " ", "container.Name: ", strings.TrimLeft(cont.Name, "/"))
			log.Debugln("MATCH: ", strings.EqualFold(strings.TrimLeft(cont.Name, "/"), name))

			if strings.EqualFold(strings.TrimLeft(cont.Name, "/"), name) {
				log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Container FOUND: ", name)
				return container, true, nil
			}
		}
	}
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Container NOT Found: ", name)
	return types.Container{}, false, nil
}
