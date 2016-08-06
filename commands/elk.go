package commands

import (
	log "github.com/Sirupsen/logrus"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldocker"
)

func cmdELK(logs bool) error {

	docker := maldocker.NewDockerClient()

	// contJSON, err := &docker StartELK(logs)
	contJSON, err := docker.StartELK(logs)
	er.CheckError(err)

	log.WithFields(log.Fields{
		// "id":   cont.ID,
		"ip": docker.GetIP(),
		// "url":      "http://" + maldocker.GetIP(),
		"username": "admin",
		"password": "admin",
		"name":     contJSON.Name,
		// "env":      config.Conf.Environment.Run,
	}).Info("ELK Container Started")

	return nil
}
