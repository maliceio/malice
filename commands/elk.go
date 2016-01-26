package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/maldocker"
)

func cmdELK(logs bool) {

	docker := maldocker.NewDockerClient()

	// cont, err := docker.StartELK(logs)
	// er.CheckError(err)

	log.WithFields(log.Fields{
		// "id":   cont.ID,
		"ip": docker.GetIP(),
		// "url":      "http://" + maldocker.GetIP(),
		"username": "admin",
		"password": "admin",
		"name":     cont.Name,
		"env":      config.Conf.Environment.Run,
	}).Info("ELK Container Started")
}
