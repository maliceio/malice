package client

import (
	"os"
	"os/exec"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"
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
		log.Debug("Creating NewClient...")
		proto, addr, basePath, err := client.ParseHost("unix:///var/run/docker.sock")
		log.Debug("Proto: ", proto, ", Addr: ", addr, ", BasePath: ", basePath, ", Error: ", err)
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		docker, err = client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
		ip = "localhost"
		port = "2375"
	case "darwin":
		log.Debug("Running inside Docker for Mac...")
		log.Debug("Creating NewClient...")
		proto, addr, basePath, err := client.ParseHost("unix:///var/run/docker.sock")
		log.Debug("Proto: ", proto, ", Addr: ", addr, ", BasePath: ", basePath, ", Error: ", err)
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		docker, err = client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
		ip = "localhost"
		port = "2375"
	case "windows":
		log.Debug("Creating NewEnvClient...")
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
	log.Debug("Docker Info...")
	if _, err = docker.Info(context.Background()); err != nil {
		log.Debug("Docker Info FAILED...")
		handleClientError(err)
	} else {
		log.WithFields(log.Fields{"ip": ip, "port": port}).Debug("Connected to docker daemon client")
	}

	log.Debug("Docker Info...")
	if _, err = docker.Info(context.Background()); err != nil {
		log.Debug("Docker Info FAILED...")
		handleClientError(err)
	} else {
		log.WithFields(log.Fields{"ip": ip, "port": port}).Debug("Connected to docker daemon client")
	}

	log.Debug("Docker Info...")
	if _, err = docker.Info(context.Background()); err != nil {
		log.Debug("Docker Info FAILED...")
		er.CheckError(err)
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

// TODO: Make this betta MUCHO betta
func handleClientError(dockerError error) {
	if dockerError != nil {
		log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Error("Unable to connect to docker client")
		switch runtime.GOOS {
		case "darwin":
			if _, err := exec.LookPath("docker-machine"); err != nil {
				log.Info("Please install docker-machine by running: ")
				log.Info(" - brew install docker-machine")
				log.Infof(" - brew install docker-machine\n\tdocker-machine create -d virtualbox %s", config.Conf.Docker.Name)
				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
			} else {
				log.Info("Please start and source the docker-machine env by running: ")
				log.Infof(" - docker-machine start %s", config.Conf.Docker.Name)
				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
			}
		case "linux":
			log.Info("Please start the docker daemon.")
		case "windows":
			if _, err := exec.LookPath("docker-machine.exe"); err != nil {
				log.Info("Please install docker-machine - https://www.docker.com/docker-toolbox")
			} else {
				log.Info("Please start and source the docker-machine env by running: ")
				log.Infof(" - docker-machine start %", config.Conf.Docker.Name)
				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
			}
		}
		// TODO Decide if I want to make docker machines or rely on user to create their own.
		// log.Info("Trying to create new docker-machine: ", "test")
		// MakeDockerMachine("test")
		os.Exit(2)
	}
}
