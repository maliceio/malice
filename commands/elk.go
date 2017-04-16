package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/malice/database/elasticsearch"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/ui"
)

func cmdELK(logs bool) error {

	docker := client.NewDockerClient()

	_, err := elasticsearch.Start(docker, logs)
	if err != nil {
		log.Error(err)
	}

	_, err = ui.Start(docker, logs)
	if err != nil {
		log.Error(err)
	}

	return nil
}
