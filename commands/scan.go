package commands

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/malice/database"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/plugins"
	"github.com/maliceio/malice/utils"
)

func cmdScan(path string, logs bool) error {
	if len(path) > 0 {
		// Check that file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatal(path + ": no such file or directory")
		}

		docker := maldocker.NewDockerClient()

		// Check that RethinkDB is running
		if _, running, _ := docker.ContainerRunning("rethink"); !running {
			log.Error("RethinkDB is NOT running, starting now...")
			rethink, err := docker.StartRethinkDB(false)
			er.CheckError(err)
			rInfo, err := docker.Client.ContainerInspect(context.Background(), rethink.ID)
			er.CheckError(err)
			er.CheckError(database.TestConnection(rInfo.NetworkSettings.IPAddress))
		}

		// Setup rethinkDB
		rInfo, err := docker.Client.ContainerInspect(context.Background(), "rethink")
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
		// Write all file data to the Database
		resp := database.WriteFileToDatabase(file)
		// os.Exit(0)
		/////////////////////////////////////////////////////////////////
		// Run all Intel Plugins on the md5 hash associated with the file
		plugins.RunIntelPlugins(docker, file.MD5, resp.GeneratedKeys[0], true)

		log.Debug("Looking for plugins that will run on: ", file.Mime)
		// Iterate over all applicable installed plugins
		plugins := plugins.GetPluginsForMime(file.Mime, true)
		log.Debug("Found these plugins: ")
		for _, plugin := range plugins {
			log.Debugf(" - %v", plugin.Name)
		}

		for _, plugin := range plugins {
			log.Debugf(">>>>> RUNNING Plugin: %s >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", plugin.Name)
			// go func() {
			// Start Plugin Container
			// TODO: don't use the default of true for --logs
			cont, err := plugin.StartPlugin(docker, file.SHA256, true)
			er.CheckError(err)

			log.WithFields(log.Fields{
				"id": cont.ID,
				"ip": docker.GetIP(),
				// "url":      "http://" + maldocker.GetIP(),
				"name": cont.Name,
				// "env":  config.Conf.Environment.Run,
			}).Debug("Plugin Container Started")

			docker.RemoveContainer(cont, false, false, false)

			// }()
			// Clean up the Plugin Container
			// TODO: I want to reuse these containers for speed eventually.

			// time.Sleep(10 * time.Millisecond)
		}
		// time.Sleep(60 * time.Second)
		log.Debug("Done with plugins.")
	} else {
		log.Error("Please supply a valid file to scan.")
	}
	return nil
}

// APIScan is an API wrapper for cmdScan
func APIScan(file string) error {
	return cmdScan(file, false)
}
