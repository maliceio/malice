package maldocker

import (
	"github.com/maliceio/malice/config"

	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

func init() {
	if config.Conf.Environment.Run == "production" {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.InfoLevel)
		// log.SetFormatter(&logstash.LogstashFormatter{Type: "malice"})
	} else {
		// Log as ASCII formatter.
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)
}

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
