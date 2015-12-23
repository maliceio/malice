package docker

import (
	"fmt"
	"os"
	"time"

	"github.com/blacktop/go-malice/config"
	"github.com/docker/docker/pkg/units"
	"github.com/fsouza/go-dockerclient"
)

var (
	endpoint = config.Conf.Malice.Docker.Endpoint
	path     = os.Getenv("DOCKER_CERT_PATH")
	ca       = fmt.Sprintf("%s/ca.pem", path)
	cert     = fmt.Sprintf("%s/cert.pem", path)
	key      = fmt.Sprintf("%s/key.pem", path)
)

// StartELK creates an ELK container from the image blacktop/elk
func StartELK() (*docker.Container, error) {
	client, _ := docker.NewTLSClient(endpoint, cert, key, ca)

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

func info() {
	var created string
	var size string

	client, _ := docker.NewTLSClient(endpoint, cert, key, ca)

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
}
