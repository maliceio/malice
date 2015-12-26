package docker

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/blacktop/go-malice/config"

	"os"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/vendor/src/github.com/docker/go-units"
	"github.com/fsouza/go-dockerclient"
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

// PingDockerClient pings docker client to see if it is up or not.
func PingDockerClient() error {
	err := client.Ping()
	if err != nil {
		log.Errorln(err)
	}
	return nil
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
					log.WithFields(log.Fields{"env": config.Conf.Malice.Environment}).Debug("Container FOUND: ", name)

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

// ImageExists returns APIImages images list and true
// if the image name exists, otherwise false.
func ImageExists(name string) (docker.APIImages, bool) {
	return ParseImages(name)
}

// ParseImages parses the images
func ParseImages(name string) (docker.APIImages, bool) {
	log.WithFields(log.Fields{
		"env": config.Conf.Malice.Environment,
	}).Debug("Searching for image: ", name)
	images := listImages()

	r := regexp.MustCompile(name)

	if len(images) != 0 {
		for _, image := range images {
			for _, tag := range image.RepoTags {
				if r.MatchString(tag) {
					log.WithFields(log.Fields{"env": config.Conf.Malice.Environment}).Debug("Image FOUND: ", name)
					return image, true
				}
			}
		}
	}

	log.WithFields(log.Fields{
		"env": config.Conf.Malice.Environment,
	}).Debug("Container NOT Found: ", name)

	return docker.APIImages{}, false
}

func listImages() []docker.APIImages {
	var images []docker.APIImages

	imageList, _ := client.ListImages(docker.ListImagesOptions{})
	for _, image := range imageList {
		images = append(images, image)
	}

	return images
}

// parseDockerEndoint returns ip and port from docker endpoint string
func parseDockerEndoint(endpoint string) (string, string, error) {

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return "", "", fmt.Errorf("Unalbe to parse endpoint: %s", endpoint)
	}

	return hostParts[0], hostParts[1], nil
}
