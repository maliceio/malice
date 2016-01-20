package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/cert"
	"github.com/samalba/dockerclient"
)

// Callback used to listen to Docker's events
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.Printf("Received event: %#v\n", *event)
}

func main() {
	// Init the client
	// docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())
	host, _ := api.Load("malice")
	url, _ := host.URL()
	tlsConfig, err := cert.ReadTLSConfig(url, host.AuthOptions())
	if err != nil {
		fmt.Printf("Unable to read TLS config: %s", err)
	}
	os.Unsetenv("DOCKER_TLS_VERIFY")
	docker, _ := dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), tlsConfig)

	// Get only running containers
	containers, err := docker.ListContainers(false, false, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range containers {
		log.Println(c.Id, c.Names)
	}

	// Inspect the first container returned
	if len(containers) > 0 {
		id := containers[0].Id
		info, _ := docker.InspectContainer(id)
		log.Println(info)
	}

	// // Build a docker image
	// // some.tar contains the build context (Dockerfile any any files it needs to add/copy)
	// dockerBuildContext, err := os.Open("some.tar")
	// defer dockerBuildContext.Close()
	// buildImageConfig := &dockerclient.BuildImage{
	// 	Context:        dockerBuildContext,
	// 	RepoName:       "your_image_name",
	// 	SuppressOutput: false,
	// }
	// reader, err := docker.BuildImage(buildImageConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var authConfig *dockerclient.AuthConfig
	err = docker.PullImage("ubuntu:14.04", authConfig)
	if err != nil {
		fmt.Printf("Unable to pull Image: %s", err)
	}

	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:       "ubuntu:14.04",
		Cmd:         []string{"bash"},
		AttachStdin: true,
		Tty:         true}
	containerID, err := docker.CreateContainer(containerConfig, "foobar", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("containerId: ", containerID)
	// Start the container
	hostConfig := &dockerclient.HostConfig{}
	err = docker.StartContainer(containerID, hostConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Stop the container (with 5 seconds timeout)
	docker.StopContainer(containerID, 5)

	// Listen to events
	docker.StartMonitorEvents(eventCallback, nil)

	// Hold the execution to look at the events coming
	time.Sleep(3600 * time.Second)
}
