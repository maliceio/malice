package client

import (
	"os"
	"os/exec"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/maliceio/malice/config"
	"golang.org/x/net/context"
)

// NOTE: https://github.com/eris-ltd/eris-cli/blob/master/perform/docker_run.go

// Docker is the Malice docker client
type Docker struct {
	Client *client.Client
	ip     string
	port   string
}

// NewDockerClient creates a new Docker Client
func NewDockerClient() *Docker {
	var docker *client.Client
	var ip, port string
	var err error

	switch os := runtime.GOOS; os {
	case "linux":
		log.Debug("Running inside Docker...")
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		docker, err = client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
		ip = "localhost"
		port = "2375"
	case "darwin":
		log.Debug("Running on Docker for Mac...")
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		docker, err = client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
		ip = "localhost"
		port = "2375"
	case "windows":
		log.Debug("Running on Docker for Windows or docker-machine on a Windows host...")
		docker, err = client.NewEnvClient()
		if err != nil {
			log.Fatal(err)
		}
		ip, port, err = parseDockerEndoint(utils.Getopt("DOCKER_HOST", config.Conf.Docker.EndPoint))
	default:
		log.Debug("Creating NewEnvClient...")
		docker, err = client.NewEnvClient()
		if err != nil {
			log.Fatal(err)
		}
		ip, port, err = parseDockerEndoint(utils.Getopt("DOCKER_HOST", config.Conf.Docker.EndPoint))
	}
	if err != nil {
		log.Fatal(err)
	}
	// Check if client can connect
	if _, err = docker.Info(context.Background()); err != nil {
		handleClientError(err)
	} else {
		log.WithFields(log.Fields{"ip": ip, "port": port}).Debug("Connected to docker daemon client")
	}

	return &Docker{
		Client: docker,
		ip:     ip,
		port:   port,
	}
}

// GetIP returns IP of docker client
func (docker *Docker) GetIP() string {
	return docker.ip
}

// TODO: Make this betta MUCHO betta
func handleClientError(dockerError error) {
	if dockerError != nil {
		log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Error("Unable to connect to docker client")
		switch runtime.GOOS {
		case "darwin":
			if _, err := os.Stat("/Applications/Docker.app"); os.IsNotExist(err) {
				log.Info("Please install Docker for Mac - https://docs.docker.com/docker-for-mac/")
				log.Info("= OR =")
				log.Info("Please install docker-machine by running: ")
				log.Info(" - brew install docker-machine")
				log.Infof(" - docker-machine create -d virtualbox %s", config.Conf.Docker.Name)
				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
			} else {
				log.Info("Please start Docker for Mac.")
				log.Info("= OR =")
				log.Info("Please start and source the docker-machine env by running: ")
				log.Infof(" - docker-machine start %s", config.Conf.Docker.Name)
				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
			}
		case "linux":
			log.Info("Please start the docker daemon. `sudo service docker start`")
		case "windows":
			if _, err := exec.LookPath("/Applications/Docker.app"); err != nil {
				log.Info("Please install Docker for Windows - https://docs.docker.com/docker-for-windows/")
				log.Info("= OR =")
				log.Info("Please install docker-toolbox - https://www.docker.com/docker-toolbox")
			} else {
				log.Info("Please start Docker for Windows.")
				log.Info("= OR =")
				log.Info("Please start and source the docker-machine env by running: ")
				log.Infof(" - docker-machine start %", config.Conf.Docker.Name)
				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
			}
		}
		os.Exit(2)
	}
}
