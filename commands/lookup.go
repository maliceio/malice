package commands

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/malice-plugins/pkgs/database/elasticsearch"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/plugins"
	"github.com/pkg/errors"
)

func cmdLookUp(hash string, logs bool) error {

	docker := client.NewDockerClient()

	elasticsearchInDocker := false
	es := elasticsearch.Database{
		Index:    utils.Getopt("MALICE_ELASTICSEARCH_INDEX", "malice"),
		Type:     utils.Getopt("MALICE_ELASTICSEARCH_TYPE", "samples"),
		URL:      utils.Getopt("MALICE_ELASTICSEARCH_URL", config.Conf.DB.URL),
		Username: utils.Getopt("MALICE_ELASTICSEARCH_USERNAME", config.Conf.DB.Username),
		Password: utils.Getopt("MALICE_ELASTICSEARCH_PASSWORD", config.Conf.DB.Password),
	}

	// This assumes you haven't set up an elasticsearch instance and that malice should create one
	if strings.EqualFold(es.URL, "http://localhost:9200") {
		elasticsearchInDocker = true
		// Check that database is running
		if _, running, _ := container.Running(docker, config.Conf.DB.Name); !running {
			log.Error("database is NOT running, starting now...")
			err := database.Start(docker, es, logs)
			if err != nil {
				return errors.Wrap(err, "failed to start to database")
			}
		}
	}

	// Initialize the malice database
	es.Init()

	if plugins.InstalledPluginsCheck(docker) {
		log.Debug("All enabled plugins are installed.")
	} else {
		// Prompt user to install all plugins?
		fmt.Println("All enabled plugins not installed would you like to install them now? (yes/no)")
		fmt.Println("[Warning] This can take a while if it is the first time you have ran Malice.")
		if utils.AskForConfirmation() {
			plugins.UpdateAllPlugins(docker)
		}
	}

	/////////////////////////////
	// Write hash to the Database
	resp, err := es.StoreHash(hash)
	if err != nil {
		return errors.Wrap(err, "cmd lookup failed to store hash")
	}

	plugins.RunIntelPlugins(docker, hash, resp.Id, true, elasticsearchInDocker)

	return nil
}

// APILookUp is an API wrapper for cmdLookUp
func APILookUp(hash string) error {
	return cmdLookUp(hash, false)
}
