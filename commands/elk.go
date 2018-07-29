package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/database/elasticsearch"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/malice/ui"
	"github.com/pkg/errors"
)

func cmdELK(logs bool) error {

	docker := client.NewDockerClient()

	if _, running, _ := container.Running(docker, config.Conf.DB.Name); !running {
		err := database.Start(docker, elasticsearch.Database{Host: config.Conf.DB.Server}, logs)
		if err != nil {
			return errors.Wrap(err, "failed to start to database")
		}
	} else {
		log.Warnf("container %s is already running", config.Conf.DB.Name)
	}

	if _, running, _ := container.Running(docker, config.Conf.UI.Name); !running {
		_, err := ui.Start(docker, logs)
		if err != nil {
			return errors.Wrap(err, "failed to start to UI")
		}
	} else {
		log.Warnf("container %s is already running", config.Conf.UI.Name)
	}

	return nil
}
