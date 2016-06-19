package commands

import (
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/maliceio/malice/plugins"
)

func cmdLookUp(hash string, logs bool) error {

	docker := maldocker.NewDockerClient()

	plugins.RunIntelPlugins(docker, hash, true)

	return nil
}
