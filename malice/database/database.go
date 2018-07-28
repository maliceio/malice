package database

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	. "github.com/malice-plugins/go-plugin-utils/database/elasticsearch"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/plugins"
	"github.com/maliceio/malice/utils"
)

// GetPluginsByCategory gets malice plugins organized by category
func GetPluginsByCategory() map[string]interface{} {
	categoryList := make(map[string]interface{})
	for _, category := range plugins.GetCategories() {
		pluginList := make(map[string]interface{})
		for _, plugin := range plugins.GetAllPluginsInCategory(category) {
			pluginList[plugin.Name] = nil
		}
		categoryList[category] = pluginList
	}

	return categoryList
}

func getElasticSearchAddr(addr string) string {
	if _, inDocker := os.LookupEnv("MALICE_IN_DOCKER"); inDocker {
		if addr != "" {
			return fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", addr))
		}
		return fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elasticsearch"))
	}
	return fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "localhost"))
}

// Start creates an Elasticsearch container from the image blacktop/elasticsearch
func Start(docker *client.Docker, logs bool) (types.ContainerJSONBase, error) {

	name := config.Conf.DB.Name
	image := config.Conf.DB.Image
	binds := []string{"malice:/usr/share/elasticsearch/data"}
	portBindings := nat.PortMap{
		"9200/tcp": {{HostIP: "0.0.0.0", HostPort: "9200"}},
	}

	if docker.Ping() {
		cont, err := container.Start(docker, nil, name, image, logs, binds, portBindings, nil, nil)
		if err != nil {
			return types.ContainerJSONBase{}, err
		}
		// Inspect newly created container to get IP assigned to it
		dbInfo, err := container.Inspect(docker, cont.ID)
		if err != nil {
			log.Error(err)
		}
		elasticAddress := getElasticSearchAddr(dbInfo.NetworkSettings.IPAddress)

		log.WithFields(log.Fields{
			// "id":   cont.ID,
			"ip":   docker.GetIP(),
			"port": config.Conf.DB.Ports,
			"name": cont.Name,
			"env":  config.Conf.Environment.Run,
		}).Info("Elasticsearch Container Started")

		// Give ELK a few seconds to start
		log.WithFields(log.Fields{
			"server":  elasticAddress,
			"timeout": config.Conf.DB.Timeout,
		}).Info("Waiting for Elasticsearch to come online.")

		ctx := context.Background()
		err = WaitForConnection(ctx, "", config.Conf.DB.Timeout)
		if err != nil {
			logOpts := types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Follow:     false,
			}

			logs, _ := docker.Client.ContainerLogs(context.Background(), cont.ID, logOpts)
			defer logs.Close()
			// Convert logs to a string
			buf := new(bytes.Buffer)
			buf.ReadFrom(logs)
			logStr := buf.String()
			log.Debug(logStr)
			// Check if elasticsearch could not start due to lack of RAM
			if strings.Contains(logStr, "There is insufficient memory for the Java Runtime Environment to continue") {
				info, err := docker.Client.Info(context.Background())
				if err != nil {
					log.Error(err)
				}
				log.Fatal("You do not have enough RAM to run elasticsearch. Elasticsearch needs at least 2GB and you have: ", units.BytesSize(float64(info.MemTotal)))
			}
			if err != nil {
				log.Error(err)
			}
		}

		return cont, err
	}
	return types.ContainerJSONBase{}, errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
}

// WaitForConnection waits for connection to Elasticsearch to be ready
func WaitForConnection(ctx context.Context, addr string, timeout int) error {

	var ready bool
	var connErr error
	secondsWaited := 0

	connCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	log.Debug("===> trying to connect to elasticsearch")
	for {
		// Try to connect to Elasticsearch
		select {
		case <-connCtx.Done():
			log.WithFields(log.Fields{"timeout": timeout}).Error("connecting to elasticsearch timed out")
			return connErr
		default:
			ready, connErr = TestConnection(addr)
			if ready {
				log.Infof("Elasticsearch came online after %d seconds", secondsWaited)
				return connErr
			}
			secondsWaited++
			time.Sleep(1 * time.Second)
		}
	}
}
