package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/libmalice/maldocker"
	"github.com/maliceio/malice/libmalice/persist"
	"github.com/maliceio/malice/plugins"
)

func init() {
	if config.Conf.Environment.Run == "production" {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.InfoLevel)
		// log.SetFormatter(&logstash.LogstashFormatter{Type: "malice"})
	} else {
		// Log as ASCII formatter.
		log.SetFormatter(&log.TextFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)
}

func cmdScan(path string, logs bool) {

	file := persist.File{
		Path: path,
	}

	// file.Mime, _ = file.GetFileMimeType()
	// file.GetFileMimeType()
	file.Init()

	log.Debug("Looking for plugins that will run on: ", file.Mime)
	// Iterate over all applicable plugins
	plugins := plugins.GetPluginsForMime(file.Mime)
	log.Debug("Found these plugins: ", plugins)
	for _, plugin := range plugins {
		cont, err := plugin.StartPlugin(logs)
		assert(err)

		log.WithFields(log.Fields{
			"id": cont.ID,
			"ip": maldocker.GetIP(),
			// "url":      "http://" + maldocker.GetIP(),
			"name": cont.Name,
			"env":  config.Conf.Environment.Run,
		}).Debug("Plugin Container Started")
	}
	log.Debug("Done with plugins.")
	// Output File Hashes
	file.PrintFileDetails()

	log.WithFields(log.Fields{
		"mime": file.Mime,
		"path": path,
		"env":  config.Conf.Environment.Run,
	}).Debug("Mime Type")
}
