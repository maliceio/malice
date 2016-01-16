package maldocker

import (
	"github.com/maliceio/malice/config"

	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

// ContainerExists returns APIContainers containers list and true
// if the container name exists, otherwise false.
func ContainerExists(client *docker.Client, name string) (*docker.APIContainers, bool, error) {
	return ParseContainers(client, name, true)
}

// ContainerRunning returns APIContainers containers list and true
// if the container name exists and is running, otherwise false.
func ContainerRunning(client *docker.Client, name string) (*docker.APIContainers, bool, error) {
	return ParseContainers(client, name, false)
}

// ParseContainers parses the containers
func ParseContainers(client *docker.Client, name string, all bool) (*docker.APIContainers, bool, error) {
	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Debug("Searching for container: ", name)
	containers, err := listContainers(client, all)
	if err != nil {
		return nil, false, err
	}

	r := regexp.MustCompile(name)

	if len(containers) != 0 {
		for _, container := range containers {
			for _, n := range container.Names {
				if r.MatchString(n) {
					log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Container FOUND: ", name)

					return &container, true, nil
				}
			}
		}
	}

	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Debug("Container NOT Found: ", name)

	return nil, false, nil
}

// listContainers returns array of APIContainers and error
func listContainers(client *docker.Client, all bool) ([]docker.APIContainers, error) {
	var containers []docker.APIContainers

	containerList, err := client.ListContainers(docker.ListContainersOptions{All: all})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for _, container := range containerList {
		containers = append(containers, container)
	}

	return containers, nil
}
