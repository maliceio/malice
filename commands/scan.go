package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/plugins"
)

func cmdScan(path string, logs bool) {

	docker := maldocker.NewDockerClient()

	file := persist.File{
		Path: path,
	}

	file.Init()
	// Output File Hashes
	log.Debug("[File]")
	file.ToMarkdownTable()
	// fmt.Println(string(file.ToJSON()))

	log.Debug("Looking for plugins that will run on: ", file.Mime)
	// Iterate over all applicable plugins
	plugins := plugins.GetPluginsForMime(file.Mime)
	log.Debug("Found these plugins: ", plugins)
	for _, plugin := range plugins {
		log.Debugf("[%s]\n", plugin.Name)
		cont, err := plugin.StartPlugin(docker, file.SHA256, logs)
		er.CheckError(err)

		log.WithFields(log.Fields{
			"id": cont.ID,
			"ip": docker.GetIP(),
			// "url":      "http://" + maldocker.GetIP(),
			"name": cont.Name,
			"env":  config.Conf.Environment.Run,
		}).Debug("Plugin Container Started")
		// Clean up the Plugin Container
		// TODO: I want to reuse these containers for speed eventually.
		err = docker.ContainerRemove(cont, false, false)
		er.CheckError(err)
	}
	log.Debug("Done with plugins.")

	// file.PrintFileDetails()

	log.WithFields(log.Fields{
		"mime": file.Mime,
		"path": path,
		"env":  config.Conf.Environment.Run,
	}).Debug("Mime Type")
}
