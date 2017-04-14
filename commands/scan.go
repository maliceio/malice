package commands

import (
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database/elasticsearch"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/plugins"
	"github.com/maliceio/malice/utils"
)

func cmdScan(path string, logs bool) error {

	ScanSample(path)

	return nil
}

// APIScan is an API wrapper for cmdScan
func APIScan(file string) error {
	return cmdScan(file, false)
}

// ScanSample scans a sample with all appropreiate malice plugins
func ScanSample(path string) {
	if len(path) > 0 {
		// Check that file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatal(path + ": no such file or directory")
		}

		docker := client.NewDockerClient()

		// Check that ElasticSearch is running
		if _, running, _ := container.Running(docker, config.Conf.DB.Name); !running {
			log.Error("Elasticsearch is NOT running, starting now...")
			_, err := elasticsearch.Start(docker, false)
			er.CheckError(err)
		}

		// Setup ElasticSearch
		dbInfo, err := container.Inspect(docker, config.Conf.DB.Name)
		er.CheckError(err)
		log.WithFields(log.Fields{
			"ip":      dbInfo.NetworkSettings.IPAddress,
			"network": dbInfo.HostConfig.NetworkMode,
			"image":   dbInfo.Config.Image,
		}).Debug("Elasticsearch is running.")

		elasticsearch.InitElasticSearch(dbInfo.NetworkSettings.IPAddress)

		if plugins.InstalledPluginsCheck(docker) {
			log.Debug("All enabled plugins are installed.")
		} else {
			// Prompt user to install all plugins?
			fmt.Println("All enabled plugins not installed would you like to install them now? (yes/no)")
			fmt.Println("[Warning] This can take a while if it is the first time you have ran Malice.")
			if utils.AskForConfirmation() {
				plugins.UpdateEnabledPlugins(docker)
			}
		}

		file := persist.File{Path: path}
		file.Init()

		// Output File Hashes
		file.ToMarkdownTable()
		// fmt.Println(string(file.ToJSON()))

		//////////////////////////////////////
		// Copy file to malice volume
		container.CopyToVolume(docker, file)

		//////////////////////////////////////
		// Write all file data to the Database
		resp := elasticsearch.WriteFileToDatabase(structs.Map(file))
		scanID := resp.Id

		/////////////////////////////////////////////////////////////////
		// Run all Intel Plugins on the md5 hash associated with the file
		plugins.RunIntelPlugins(docker, file.SHA1, scanID, true)

		// Get file's mime type
		mimeType, err := persist.GetMimeType(docker, file.SHA256)
		er.CheckError(err)

		log.Debug("Looking for plugins that will run on: ", mimeType)
		// Iterate over all applicable installed plugins
		plugins := plugins.GetPluginsForMime(mimeType, true)
		log.Debug("Found these plugins: ")
		for _, plugin := range plugins {
			log.Debugf(" - %v", plugin.Name)
		}

		var wg sync.WaitGroup
		wg.Add(len(plugins))

		for _, plugin := range plugins {
			log.Debugf(">>>>> RUNNING Plugin: %s >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", plugin.Name)
			// Start Plugin Container
			// TODO: don't use the default of true for --logs
			go plugin.StartPlugin(docker, file.SHA256, scanID, true, &wg)
		}

		wg.Wait() // this waits for the counter to be 0
		log.Debug("Done with plugins.")
	} else {
		log.Error("Please supply a valid file to scan.")
	}
}
