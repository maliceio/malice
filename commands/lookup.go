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

// APILookUp is an API wrapper for cmdLookUp
func APILookUp(hash string) error {
	return cmdLookUp(hash, false)
}
