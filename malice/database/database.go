package database

import (
	"bytes"
	"context"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	"github.com/malice-plugins/pkgs/database/elasticsearch"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/plugins"
	"github.com/pkg/errors"
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

// Start creates an Elasticsearch container from the image blacktop/elasticsearch
func Start(docker *client.Docker, es elasticsearch.Database, logs bool) error {

	name := config.Conf.DB.Name
	image := config.Conf.DB.Image
	binds := []string{"malice:/usr/share/elasticsearch/data"}
	portBindings := nat.PortMap{
		"9200/tcp": {{HostIP: "0.0.0.0", HostPort: "9200"}},
	}

	if docker.Ping() {
		esContainer, err := container.Start(docker, nil, name, image, logs, binds, portBindings, nil, nil)
		if err != nil {
			return errors.Wrap(err, "failed to start docker container")
		}

		// Inspect newly created container to get IP assigned to it
		dbInfo, err := container.Inspect(docker, esContainer.ID)
		if err != nil {
			return errors.Wrap(err, "failed to inspect container")
		}

		log.WithFields(log.Fields{
			// "id":   cont.ID,
			"docker_ip":   docker.GetIP(),
			"assigned_ip": dbInfo.NetworkSettings.IPAddress,
			"port":        config.Conf.DB.Ports,
			"name":        esContainer.Name,
			"runtime_env": config.Conf.Environment.Run,
		}).Info("elasticsearch container started")

		// Wait for Elasticsearch to start (takes ~10-20 secs)
		err = es.WaitForConnection(context.Background(), config.Conf.DB.Timeout)
		if err != nil {
			logOpts := types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Follow:     false,
			}

			logs, _ := docker.Client.ContainerLogs(context.Background(), esContainer.ID, logOpts)
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
					return err
				}
				log.Fatal("You do not have enough RAM to run elasticsearch. Elasticsearch needs at least 2GB and you have: ", units.BytesSize(float64(info.MemTotal)))
			}
			if err != nil {
				return err
			}
		}

		return nil
	}
	return errors.New("cannot connect to the Docker daemon")
}

// Setup ElasticSearch
//dbInfo, err := container.Inspect(docker, config.Conf.DB.Name)
//er.CheckError(err)
//log.WithFields(log.Fields{
//"ip":      dbInfo.NetworkSettings.IPAddress,
//"network": dbInfo.HostConfig.NetworkMode,
//"image":   dbInfo.Config.Image,
//}).Debug("Elasticsearch is running.")
//

//db.Init()
