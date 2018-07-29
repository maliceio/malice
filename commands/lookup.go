package commands

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/database/elasticsearch"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/plugins"
	"github.com/maliceio/malice/utils"
	"github.com/pkg/errors"
)

func cmdLookUp(hash string, logs bool) error {

	docker := client.NewDockerClient()
	es := elasticsearch.Database{
		Host:  config.Conf.DB.Server,
		Index: "malice",
		Type:  "samples",
	}

	// Check that database is running
	if _, running, _ := container.Running(docker, config.Conf.DB.Name); !running {
		log.Error("database is NOT running, starting now...")
		err := database.Start(docker, es, logs)
		if err != nil {
			return errors.Wrap(err, "failed to start to database")
		}
		// Initialize the malice database
		es.Init()
	}

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

	plugins.RunIntelPlugins(docker, hash, resp.Id, true)

	return nil
}

// APILookUp is an API wrapper for cmdLookUp
func APILookUp(hash string) error {
	return cmdLookUp(hash, false)
}
