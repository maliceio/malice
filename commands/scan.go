package commands

import (
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/malice/database"
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

		// Check that RethinkDB is running
		if _, running, _ := container.Running(docker, "rethink"); !running {
			log.Error("RethinkDB is NOT running, starting now...")
			rethink, err := container.StartRethinkDB(docker, false)
			er.CheckError(err)
			rInfo, err := container.Inspect(docker, rethink.ID)
			er.CheckError(err)
			er.CheckError(database.TestConnection(rInfo.NetworkSettings.IPAddress))
		}

		// Setup rethinkDB
		rInfo, err := container.Inspect(docker, "rethink")
		er.CheckError(err)
		er.CheckError(database.TestConnection(rInfo.NetworkSettings.IPAddress))
		database.InitRethinkDB()

		if plugins.InstalledPluginsCheck(docker) {
			log.Debug("All enabled plugins are installed.")
		} else {
			// Prompt user to install all plugins?
			fmt.Println("All enabled plugins not installed would you like to install them now? (yes/no)")
			fmt.Println("[Warning] This can take a while if it is the first time you have ran Malice.")
			if util.AskForConfirmation() {
				plugins.UpdateAllPlugins(docker)
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
		resp := database.WriteFileToDatabase(file)
		scanID := resp.GeneratedKeys[0]

		/////////////////////////////////////////////////////////////////
		// Run all Intel Plugins on the md5 hash associated with the file
		plugins.RunIntelPlugins(docker, file.MD5, scanID, true)

		log.Debug("Looking for plugins that will run on: ", file.Mime)
		// Iterate over all applicable installed plugins
		plugins := plugins.GetPluginsForMime(file.Mime, true)
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
