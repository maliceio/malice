package maldocker

import (
	"os"
	"os/exec"
	"runtime"

	log "github.com/Sirupsen/logrus"
	docker "github.com/docker/engine-api/client"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"
)

// NOTE: https://github.com/eris-ltd/eris-cli/blob/master/perform/docker_run.go

// Docker is the Malice docker client
type Docker struct {
	Client *docker.Client
	ip     string
	port   string
}

// NewDockerClient creates a new Docker Client
func NewDockerClient() *Docker {
	var client *docker.Client
	var ip, port string
	var err error
	// create docker client base on OS
	switch runtime.GOOS {
	case "darwin", "windows":
		ip, port, err = parseDockerEndoint(os.Getenv("DOCKER_HOST"))
		er.CheckError(err)
		client, err = docker.NewEnvClient()
		handleClientError(err)
	case "linux":
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		client, err = docker.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
		handleClientError(err)
	}
	return &Docker{
		Client: client,
		ip:     ip,
		port:   port,
	}
}

// GetIP returns IP of docker client
func (client *Docker) GetIP() string {
	return client.ip
}

// TODO: Make this betta MUCHO betta
func handleClientError(dockerError error) {
	if dockerError != nil {
		log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Error("Unable to connect to docker client")
		switch runtime.GOOS {
		case "darwin":
			if _, err := exec.LookPath("docker-machine"); err != nil {
				log.Infof("Please install docker-machine by running: \n\tbrew install docker-machine\n\tdocker-machine create -d virtualbox %s\n\teval $(docker-machine env %s)\n", config.Conf.Docker.Name, config.Conf.Docker.Name)
			} else {
				log.Infof("Please start and source the docker-machine env by running: \n\tdocker-machine start %s\n\teval $(docker-machine env %s)\n", config.Conf.Docker.Name, config.Conf.Docker.Name)
			}
		case "linux":
			log.Info("Please start the docker daemon.")
		case "windows":
			if _, err := exec.LookPath("docker-machine.exe"); err != nil {
				log.Info("Please install docker-machine - https://www.docker.com/docker-toolbox")
			} else {
				log.Infof("Please start and source the docker-machine env by running: \n\tdocker-machine start %s\n\teval $(docker-machine env %s)\n", config.Conf.Docker.Name, config.Conf.Docker.Name)
			}
		}
		// TODO Decide if I want to make docker machines or rely on use to create their own.
		// log.Info("Trying to create new docker-machine: ", "test")
		// MakeDockerMachine("test")
		os.Exit(2)
	}
}

// LogContainer tails container logs to terminal
// func LogContainer(cont *docker.Container) {
//
// 	opts := docker.LogsOptions{
// 		Container:    cont.ID,
// 		OutputStream: os.Stdout,
// 		ErrorStream:  os.Stderr,
// 		Follow:       true,
// 		Stdout:       true,
// 		Stderr:       true,
// 		// Since:        0,
// 		Timestamps: false,
// 		// Tail:         false,
// 		RawTerminal: false, // Usually true when the container contains a TTY.
// 	}
//
// 	er.CheckError(client.Logs(opts))
// }

// StartELK creates an ELK container from the image blacktop/elk
// func StartELK(logs bool) (cont *docker.Container, err error) {
//
// 	er.CheckError(PingDockerClient(client))
//
// 	_, exists, err := ContainerExists(client, "elk")
//
// 	if exists {
// 		log.WithFields(log.Fields{
// 			"exisits": exists,
// 			// "id":      elkContainer.ID,
// 			"env": config.Conf.Environment.Run,
// 			"url": "http://" + ip,
// 		}).Info("ELK is already running...")
// 		os.Exit(0)
// 	}
//
// 	_, exists, err = ImageExists(client, "blacktop/elk")
// 	if exists {
// 		log.WithFields(log.Fields{
// 			"exisits": exists,
// 			// "id":      elkContainer.ID,
// 			"env": config.Conf.Environment.Run,
// 		}).Info("Image `blacktop/elk` already pulled.")
// 	} else {
// 		log.WithFields(log.Fields{
// 			"exisits": exists,
// 			"env":     config.Conf.Environment.Run}).Info("Pulling Image `blacktop/elk`")
//
// 		er.CheckError(PullImage("blacktop/elk", "latest"))
// 	}
// 	createContConf := docker.Config{
// 		Image: "blacktop/elk",
// 	}
//
// 	portBindings := map[docker.Port][]docker.PortBinding{
// 		"80/tcp":   {{HostIP: "0.0.0.0", HostPort: "80"}},
// 		"9200/tcp": {{HostIP: "0.0.0.0", HostPort: "9200"}},
// 	}
//
// 	createContHostConfig := docker.HostConfig{
// 		// Binds:           []string{"/var/run:/var/run:rw", "/sys:/sys:ro", "/var/lib/docker:/var/lib/docker:ro"},
// 		PortBindings: portBindings,
// 		// PublishAllPorts: true,
// 		Privileged: false,
// 	}
//
// 	createContOps := docker.CreateContainerOptions{
// 		Name:       "elk",
// 		Config:     &createContConf,
// 		HostConfig: &createContHostConfig,
// 	}
//
// 	cont, err = client.CreateContainer(createContOps)
// 	if err != nil {
// 		log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Errorf("CreateContainer error = %s\n", err)
// 	}
//
// 	err = client.StartContainer(cont.ID, nil)
// 	if err != nil {
// 		log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Errorf("StartContainer error = %s\n", err)
// 	}
//
// 	if logs {
// 		LogContainer(cont)
// 	}
//
// 	return
// }
