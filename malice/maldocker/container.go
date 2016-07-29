package maldocker

import (
	"errors"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"

	log "github.com/Sirupsen/logrus"
)

func (client *Docker) checkContainerRequirements(containerName, image string) {
	// Check for existance of malice network
	if _, exists, _ := client.NetworkExists("malice"); !exists {
		log.WithFields(log.Fields{
			"network": "malice",
			"exisits": exists,
			"env":     config.Conf.Environment.Run,
		}).Error("Network malice does not exist, creating now...")
		_, err := client.CreateNetwork("malice")
		er.CheckError(err)
	}
	// Check for existance of malice volume
	if _, exists, _ := client.VolumeExists("malice"); !exists {
		log.Debug("Volume malice not found.")
		_, err := client.CreateVolume("malice")
		er.CheckError(err)
	}
	log.Debug("Volume malice found.")
	// Check that the container isn't already running
	if _, exists, _ := client.ContainerExists(containerName); exists {
		log.WithFields(log.Fields{
			"exisits": exists,
			"name":    containerName,
			"env":     config.Conf.Environment.Run,
		}).Error("Container is already running...")
		os.Exit(0)
	}
	// Check that we have already pulled the image
	if _, exists, _ := client.ImageExists(image); exists {
		log.WithFields(log.Fields{
			"exisits": exists,
			"env":     config.Conf.Environment.Run,
		}).Debugf("Image `%s` already pulled.", image)
	} else {
		log.WithFields(log.Fields{
			"exisits": exists,
			"env":     config.Conf.Environment.Run}).Debugf("Pulling Image `%s`", image)
		client.PullImage(image, "latest")
	}
}

// StartContainer starts a malice docker container
func (client *Docker) StartContainer(
	cmd strslice.StrSlice,
	name string,
	image string,
	logs bool,
	binds []string,
	portBindings nat.PortMap,
	links []string,
	env []string,
) (types.ContainerJSONBase, error) {

	if client.Ping() {
		// Check that all requirements for the container to run are ready
		client.checkContainerRequirements(name, image)

		ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Docker.Timeout*time.Second)
		defer cancel()

		createContConf := &container.Config{
			Image: image,
			Cmd:   cmd,
			Env:   env,
			// Env:   []string{"MALICE_VT_API=" + os.Getenv("MALICE_VT_API")},
		}
		hostConfig := &container.HostConfig{
			// Binds:      []string{maldirs.GetSampledsDir() + ":/malware:ro"},
			// Binds:      []string{"malice:/malware:ro"},
			Binds: binds,
			// NetworkMode:  "malice",
			PortBindings: portBindings,
			Links:        links,
			Privileged:   false,
		}
		networkingConfig := &network.NetworkingConfig{}

		contResponse, err := client.Client.ContainerCreate(ctx, createContConf, hostConfig, networkingConfig, name)
		if err != nil {
			log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Errorf("CreateContainer error = %s\n", err)
		}

		err = client.Client.ContainerStart(ctx, contResponse.ID, types.ContainerStartOptions{})
		if err != nil {
			log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Errorf("StartContainer error = %s\n", err)
		}

		if logs {
			client.LogContainer(contResponse.ID)
		}

		contJSON, err := client.ContainerInspect(contResponse.ID)
		return contJSON, err
	}
	return types.ContainerJSONBase{}, errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
}

// RemoveContainer removes the `cont` container unforcedly.
// If volumes is true, the associated volumes are removed with container.
// If links is true, the associated links are removed with container.
// If force is true, the container will be destroyed with extreme prejudice.
func (client *Docker) RemoveContainer(cont types.ContainerJSONBase, volumes bool, links bool, force bool) {
	// check if container exists
	if plugin, exists, _ := client.ContainerExists(cont.Name); exists {
		log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Removing Plugin container: ", cont.Name)
		er.CheckError(client.Client.ContainerRemove(context.Background(), plugin.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			// RemoveLinks:   links,
			Force: true,
		}))
	} else {
		// container not found
		log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Error("Plugin container does not exist. Cannot remove.")
	}
}

// LogContainer tails container logs to terminal
func (client *Docker) LogContainer(contID string) {

	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Docker.Timeout*time.Second)
	defer cancel()

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		// Since       string
		// Timestamps  bool
		Follow: true,
		// Tail        string
	}

	logs, err := client.Client.ContainerLogs(ctx, contID, options)
	defer logs.Close()
	er.CheckError(err)

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, logs)
	er.CheckError(err)
}

// ContainerInspect returns types.ContainerJSON from Container ID
// if the container name exists, otherwise false.
func (client *Docker) ContainerInspect(id string) (types.ContainerJSONBase, error) {
	contJSON, err := client.Client.ContainerInspect(context.Background(), id)
	return *contJSON.ContainerJSONBase, err
}

// ContainerExists returns APIContainers containers list and true
// if the container name exists, otherwise false.
func (client *Docker) ContainerExists(name string) (types.Container, bool, error) {
	return client.ParseContainers(strings.TrimLeft(name, "/"), true)
}

// ContainerRunning returns APIContainers containers list and true
// if the container name exists and is running, otherwise false.
func (client *Docker) ContainerRunning(name string) (types.Container, bool, error) {
	return client.ParseContainers(strings.TrimLeft(name, "/"), false)
}

// ParseContainers parses the containers
func (client *Docker) ParseContainers(name string, all bool) (types.Container, bool, error) {
	// list containers
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Searching for container: ", name)
	containers, err := client.listContainers(all)
	if err != nil {
		return types.Container{}, false, err
	}
	// locate docker container that matches name
	if len(containers) != 0 {
		for _, container := range containers {

			cont, err := client.ContainerInspect(container.ID)
			er.CheckError(err)

			log.Debugln("name: ", name, " ", "container.Name: ", strings.TrimLeft(cont.Name, "/"))
			log.Debugln("MATCH: ", strings.EqualFold(strings.TrimLeft(cont.Name, "/"), name))

			if strings.EqualFold(strings.TrimLeft(cont.Name, "/"), name) {
				log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Container FOUND: ", name)
				return container, true, nil
			}
		}
	}
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Container NOT Found: ", name)
	return types.Container{}, false, nil
}

// listContainers returns array of types.Containers and error
func (client *Docker) listContainers(all bool) ([]types.Container, error) {
	options := types.ContainerListOptions{All: all}
	containers, err := client.Client.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// StartELK creates an ELK container from the image blacktop/elk
func (client *Docker) StartELK(logs bool) (types.ContainerJSONBase, error) {

	name := "elk"
	image := "blacktop/elk"
	binds := []string{"malice:/usr/share/elasticsearch/data"}
	portBindings := nat.PortMap{
		"80/tcp":   {{HostIP: "0.0.0.0", HostPort: "80"}},
		"9200/tcp": {{HostIP: "0.0.0.0", HostPort: "9200"}},
	}

	if client.Ping() {
		cont, err := client.StartContainer(nil, name, image, logs, binds, portBindings, nil, nil)
		return cont, err
	}
	return types.ContainerJSONBase{}, errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
}

// StartRethinkDB creates an RethinkDB container from the image rethinkdb
func (client *Docker) StartRethinkDB(logs bool) (types.ContainerJSONBase, error) {

	name := "rethink"
	image := "rethinkdb"
	binds := []string{"malice:/data"}
	portBindings := nat.PortMap{
		"8080/tcp":  {{HostIP: "0.0.0.0", HostPort: "8081"}},
		"28015/tcp": {{HostIP: "0.0.0.0", HostPort: "28015"}},
	}

	if client.Ping() {
		cont, err := client.StartContainer(nil, name, image, logs, binds, portBindings, nil, nil)
		// er.CheckError(err)
		// if network, exists, _ := client.NetworkExists("malice"); exists {
		// 	err := client.ConnectNetwork(network, cont)
		// 	er.CheckError(err)
		// }

		// Give rethinkDB a few seconds to start
		time.Sleep(2 * time.Second)
		log.Info("sleeping for 2 seconds to let rethinkDB start")
		return cont, err
	}
	return types.ContainerJSONBase{}, errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
}
