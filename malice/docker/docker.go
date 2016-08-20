package docker

// import (
// 	"os"
// 	"os/exec"
// 	"runtime"

// 	"golang.org/x/net/context"

// 	log "github.com/Sirupsen/logrus"
// 	docker "github.com/docker/engine-api/client"
// 	"github.com/maliceio/malice/config"
// )

// // NOTE: https://github.com/eris-ltd/eris-cli/blob/master/perform/docker_run.go

// // Docker is the Malice docker client
// type Docker struct {
// 	Client *docker.Client
// 	ip     string
// 	port   string
// }

// // NewDockerClient creates a new Docker Client
// func NewDockerClient() *Docker {
// 	var client *docker.Client
// 	var ip, port string
// 	var err error

// 	client, err = docker.NewEnvClient()

// 	// Check if client can connect
// 	if _, err = client.Info(context.Background()); err != nil {
// 		// If failed to connect try to create docker client via socket
// 		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
// 		client, err = docker.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// Check if client can connect
// 		if _, err = client.Info(context.Background()); err != nil {
// 			handleClientError(err)
// 		} else {
// 			ip = "localhost"
// 			port = "2375"
// 			log.WithFields(log.Fields{"ip": ip, "port": port}).Debug("Connected to docker daemon native client")
// 		}
// 	} else {
// 		ip, port, err = parseDockerEndoint(os.Getenv("DOCKER_HOST"))
// 		if err != nil {
// 			log.Error(err)
// 		}
// 		log.WithFields(log.Fields{"ip": ip, "port": port}).Debug("Connected to docker daemon with docker-machine")
// 	}

// 	return &Docker{
// 		Client: client,
// 		ip:     ip,
// 		port:   port,
// 	}
// }

// // GetIP returns IP of docker client
// func (client *Docker) GetIP() string {
// 	return client.ip
// }

// // TODO: Make this betta MUCHO betta
// func handleClientError(dockerError error) {
// 	if dockerError != nil {
// 		log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Error("Unable to connect to docker client")
// 		switch runtime.GOOS {
// 		case "darwin":
// 			if _, err := exec.LookPath("docker-machine"); err != nil {
// 				log.Info("Please install docker-machine by running: ")
// 				log.Info(" - brew install docker-machine")
// 				log.Infof(" - brew install docker-machine\n\tdocker-machine create -d virtualbox %s", config.Conf.Docker.Name)
// 				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
// 			} else {
// 				log.Info("Please start and source the docker-machine env by running: ")
// 				log.Infof(" - docker-machine start %s", config.Conf.Docker.Name)
// 				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
// 			}
// 		case "linux":
// 			log.Info("Please start the docker daemon.")
// 		case "windows":
// 			if _, err := exec.LookPath("docker-machine.exe"); err != nil {
// 				log.Info("Please install docker-machine - https://www.docker.com/docker-toolbox")
// 			} else {
// 				log.Info("Please start and source the docker-machine env by running: ")
// 				log.Infof(" - docker-machine start %", config.Conf.Docker.Name)
// 				log.Infof(" - eval $(docker-machine env %s)", config.Conf.Docker.Name)
// 			}
// 		}
// 		// TODO Decide if I want to make docker machines or rely on user to create their own.
// 		// log.Info("Trying to create new docker-machine: ", "test")
// 		// MakeDockerMachine("test")
// 		os.Exit(2)
// 	}
// }
