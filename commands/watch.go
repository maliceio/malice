package commands

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/plugins"
	util "github.com/maliceio/malice/utils"
	"golang.org/x/net/context"
)

func cmdWatch(folderName string, logs bool) error {

	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Info("Malice watching folder: ", folderName)

	info, err := os.Stat(folderName)

	// Check that folder exists
	if os.IsNotExist(err) {
		log.Error("error: folder does not exist.")
		return nil
	}
	// Check that path is a folder and not a file
	if info.IsDir() {
		log.Error("error: path is not a folder")
		return nil
	}

	NewWatcher(folderName)

	return nil
}

// NewWatcher creates a new watcher for the user supplied folder
func NewWatcher(folder string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("modified file:", event.Name)
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

					file := persist.File{Path: event.Name}
					file.Init()
					// Output File Hashes
					file.ToMarkdownTable()
					// fmt.Println(string(file.ToJSON()))

					//////////////////////////////////////
					// Copy file to malice volume
					docker.CopyToVolume(file)
					//////////////////////////////////////
					// Write all file data to the Database
					resp := database.WriteFileToDatabase(file)

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
							"env":  config.Conf.Environment.Run,
						}).Debug("Plugin Container Started")

						docker.RemoveContainer(cont, false, false, false)
					}
				}
			case err := <-watcher.Errors:
				log.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(folder)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
