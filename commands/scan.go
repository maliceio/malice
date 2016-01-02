package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/libmalice/maldocker"
	"github.com/maliceio/malice/utils"
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

	util.GetFileMimeType(path)

	log.WithFields(log.Fields{
		// "id":   cont.ID,
		"ip": maldocker.GetIP(),
		// "url":      "http://" + maldocker.GetIP(),
		"path": path,
		"env":  config.Conf.Environment.Run,
	}).Info("ELK Container Started")
}
