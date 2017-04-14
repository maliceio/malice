package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database/elasticsearch"
	"github.com/maliceio/malice/malice/docker/client"
	er "github.com/maliceio/malice/malice/errors"
)

func cmdELK(logs bool) error {

	docker := client.NewDockerClient()

	contJSON, err := elasticsearch.Start(docker, logs)
	er.CheckError(err)

	log.WithFields(log.Fields{
		// "id":   cont.ID,
		"ip": docker.GetIP(),
		// "url":      "http://" + docker.GetIP(),
		"username": "admin",
		"password": "admin",
		"name":     contJSON.Name,
		"env":      config.Conf.Environment.Run,
	}).Info("Elasticsearch Container Started")

	return nil
}
