package docker

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/blacktop/go-malice/config"

	"os"
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
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)
}

// PingDockerClient pings docker client to see if it is up or not.
func PingDockerClient(client *docker.Client) error {
	err := client.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//Info prints out list of docker images and containers
func Info(client *docker.Client) (err error) {
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

// parseDockerEndoint returns ip and port from docker endpoint string
func parseDockerEndoint(endpoint string) (string, string, error) {

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return "", "", fmt.Errorf("Unable to parse endpoint: %s", endpoint)
	}

	return hostParts[0], hostParts[1], nil
}
