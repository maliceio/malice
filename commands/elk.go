package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/go-malice/config"
	"github.com/blacktop/go-malice/docker"
)

func init() {
	if config.Conf.Malice.Environment == "production" {
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

func cmdELK() error {

	cont, err := docker.StartELK()
	if err != nil {
		log.Errorf("StartELK error = %s\n", err)
	}

	log.WithFields(log.Fields{
		// "id":   cont.ID,
		"ip": docker.GetIP(),
		// "url":      "http://" + docker.GetIP(),
		"username": "admin",
		"password": "admin",
		"name":     cont.Name,
		"env":      config.Conf.Malice.Environment,
	}).Info("ELK Container Started")

	return nil
}
