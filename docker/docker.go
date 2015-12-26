package docker

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/go-malice/config"
	"github.com/fsouza/go-dockerclient"
)

// NOTE: https://github.com/eris-ltd/eris-cli/blob/master/perform/docker_run.go

var (
	endpoint    = os.Getenv("DOCKER_HOST")
	ip          string
	port        string
	path        = os.Getenv("DOCKER_CERT_PATH")
	ca          = fmt.Sprintf("%s/ca.pem", path)
	cert        = fmt.Sprintf("%s/cert.pem", path)
	key         = fmt.Sprintf("%s/key.pem", path)
	client      *docker.Client
	clientError error
)

func init() {
	var err error
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
	// fmt.Println(os.GOOS)

	if endpoint != "" {
		endpoint = config.Conf.Malice.Docker.Endpoint
	}

	ip, port, err = parseDockerEndoint(endpoint)
	if err != nil {
		log.Error(err)
	}

	client, clientError = docker.NewTLSClient(endpoint, cert, key, ca)
	// Make sure we can connect to the docker client
	if clientError != nil {
		log.WithFields(log.Fields{
			"env":      config.Conf.Malice.Environment,
			"endpoint": endpoint,
		}).Error("Unable to connect to docker client")
		os.Exit(2)
	}
}

//GetIP returns IP of docker client
func GetIP() string {
	return ip
}

// StartELK creates an ELK container from the image blacktop/elk
func StartELK() (*docker.Container, error) {
	// client, _ := docker.NewTLSClient(endpoint, cert, key, ca)
	err := PingDockerClient()
	if err != nil {
		return nil, err
	}
	_, exists := ContainerExists("elk")

	if exists {
		log.WithFields(log.Fields{
			"exisits": exists,
			// "id":      elkContainer.ID,
			"env": config.Conf.Malice.Environment,
			"url": "http://" + ip,
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
