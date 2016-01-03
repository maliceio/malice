package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/libmalice/persist"
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

	file.PrintFileDetails()

	log.WithFields(log.Fields{
		"mime": file.Mime,
		"path": path,
		"env":  config.Conf.Environment.Run,
	}).Debug("Mime Type")
}
