package docker

import (
	"fmt"
	"os"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/go-malice/config"
	"github.com/docker/docker/vendor/src/github.com/docker/go-units"
	"github.com/fsouza/go-dockerclient"
)

// NOTE: https://github.com/eris-ltd/eris-cli/blob/master/perform/docker_run.go

var (
	endpoint  = config.Conf.Malice.Docker.Endpoint
	path      = os.Getenv("DOCKER_CERT_PATH")
	ca        = fmt.Sprintf("%s/ca.pem", path)
	cert      = fmt.Sprintf("%s/cert.pem", path)
	key       = fmt.Sprintf("%s/key.pem", path)
	client, _ = docker.NewTLSClient(endpoint, cert, key, ca)
)

func init() {
	if config.Conf.Malice.Environment == "production" {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.InfoLevel)
		// log.SetFormatter(&logstash.LogstashFormatter{Type: "malice"})
	} else {
		// Log as ASCII formatter.
		log.SetFormatter(&log.TextFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)
}

// StartELK creates an ELK container from the image blacktop/elk
func StartELK() (*docker.Container, error) {
	// client, _ := docker.NewTLSClient(endpoint, cert, key, ca)
	elkContainer, exists := ContainerExists("elk")

	if exists {
		log.WithFields(log.Fields{
			"exisits": exists,
			"id":      elkContainer.ID,
			"env":     config.Conf.Malice.Environment,
		}).Info("ELK is already running...")
		os.Exit(0)
	}

	createContConf := docker.Config{
		Image: "blacktop/elk",
	}

	portBindings := map[docker.Port][]docker.PortBinding{
		"80/tcp":   {{HostIP: "0.0.0.0", HostPort: "80"}},
		"9200/tcp": {{HostIP: "0.0.0.0", HostPort: "9200"}},
	}

	createContHostConfig := docker.HostConfig{
		// Binds:           []string{"/var/run:/var/run:rw", "/sys:/sys:ro", "/var/lib/docker:/var/lib/docker:ro"},
		PortBindings: portBindings,
		// PublishAllPorts: true,
		Privileged: false,
	}

	createContOps := docker.CreateContainerOptions{
		Name:       "elk",
		Config:     &createContConf,
		HostConfig: &createContHostConfig,
	}

	cont, err := client.CreateContainer(createContOps)
	if err != nil {
		fmt.Printf("create error = %s\n", err)
	}

	err = client.StartContainer(cont.ID, nil)
	if err != nil {
		fmt.Printf("start error = %s\n", err)
	}

	return cont, err
}

//Info prints out list of docker images and containers
func Info() (err error) {
	var created string
	var size string
	// var err = nil

	// client, _ := docker.NewTLSClient(endpoint, cert, key, ca)

	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	fmt.Println("Listing All Images=================================")
	for _, img := range imgs {
		// fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags[0])
		created = units.HumanDuration(time.Now().UTC().Sub(time.Unix(img.Created, 0))) + " ago"
		size = units.HumanSize(float64(img.Size))
		fmt.Println("Created: ", created)
		fmt.Println("Size: ", size)
		// fmt.Println("VirtualSize: ", img.VirtualSize)
		// fmt.Println("ParentId: ", img.ParentID)
	}
	containers, _ := client.ListContainers(docker.ListContainersOptions{All: true})
	fmt.Println("Listing All Containers==========================================")
	for _, container := range containers {
		// fmt.Println("ID: ", container.ID)
		fmt.Println("Image: ", container.Image)
		fmt.Println("Command: ", container.Command)
		created = units.HumanDuration(time.Now().UTC().Sub(time.Unix(container.Created, 0))) + " ago"
		fmt.Println("Created: ", created)
		fmt.Println("Status: ", container.Status)
		fmt.Println("Ports: ", container.Ports)
		// fmt.Println("Created: ", container.SizeRootFs)
		// fmt.Println("Created: ", container.SizeRw)
	}
	return err
}

// ContainerExists returns APIContainers containers list and true
// if the container name exists, otherwise false.
func ContainerExists(name string) (docker.APIContainers, bool) {
	return ParseContainers(name, true)
}

// ContainerRunning returns APIContainers containers list and true
// if the container name exists and is running, otherwise false.
func ContainerRunning(name string) (docker.APIContainers, bool) {
	return ParseContainers(name, false)
}

// ParseContainers parses the containers
func ParseContainers(name string, all bool) (docker.APIContainers, bool) {
	log.WithFields(log.Fields{
		"env": config.Conf.Malice.Environment,
	}).Debug("Searching for container: ", name)
	containers := listContainers(all)

	r := regexp.MustCompile(name)

	if len(containers) != 0 {
		for _, container := range containers {
			for _, n := range container.Names {
				if r.MatchString(n) {
					log.WithFields(log.Fields{
						"env": config.Conf.Malice.Environment,
					}).Debug("Container FOUND: ", name)

					return container, true
				}
			}
		}
	}

	log.WithFields(log.Fields{
		"env": config.Conf.Malice.Environment,
	}).Debug("Container NOT Found: ", name)

	return docker.APIContainers{}, false
}

func listContainers(all bool) []docker.APIContainers {
	var containers []docker.APIContainers

	containerList, _ := client.ListContainers(docker.ListContainersOptions{All: all})
	for _, container := range containerList {
		containers = append(containers, container)
	}

	return containers
}
