package commands

import (
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/maliceio/malice/plugins"
)

func cmdLookUp(hash string, logs bool) {

	docker := maldocker.NewDockerClient()

	plugins.RunIntelPlugins(docker, hash, true)
}
