package client

import (
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	er "github.com/maliceio/malice/malice/errors"
	"golang.org/x/net/context"
)

// Ping pings docker client to see if it is up or not by checking Info.
func (docker *Docker) Ping() bool {
	// ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Docker.Timeout*time.Second)
	// defer cancel()

	_, err := docker.Client.Info(context.Background())
	if err != nil {
		er.CheckError(err)
		return false
	}
	return true
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

//Info prints out list of docker images and containers
// func (docker *Docker) Info() (err error) {
// 	var created string
// 	var size string

// 	imgs, _ := image.List(docker, "", false)
// 	fmt.Println("Listing All Images=================================")
// 	for _, img := range imgs {
// 		// fmt.Println("ID: ", img.ID)
// 		fmt.Println("RepoTags: ", img.RepoTags[0])
// 		created = units.HumanDuration(time.Now().UTC().Sub(time.Unix(img.Created, 0))) + " ago"
// 		size = units.HumanSize(float64(img.Size))
// 		fmt.Println("Created: ", created)
// 		fmt.Println("Size: ", size)
// 		// fmt.Println("VirtualSize: ", img.VirtualSize)
// 		// fmt.Println("ParentId: ", img.ParentID)
// 	}
// 	containers, _ := container.List(docker, true)
// 	fmt.Println("Listing All Containers==========================================")
// 	for _, container := range containers {
// 		// fmt.Println("ID: ", container.ID)
// 		fmt.Println("Image: ", container.Image)
// 		fmt.Println("Command: ", container.Command)
// 		created = units.HumanDuration(time.Now().UTC().Sub(time.Unix(container.Created, 0))) + " ago"
// 		fmt.Println("Created: ", created)
// 		fmt.Println("Status: ", container.Status)
// 		fmt.Println("Ports: ", container.Ports)
// 		// fmt.Println("Created: ", container.SizeRootFs)
// 		// fmt.Println("Created: ", container.SizeRw)
// 	}
// 	return err
// }
