package commands

import (
	"fmt"

	"github.com/docker/machine/libmachine/log"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/plugins"
	"github.com/maliceio/malice/utils"
)

func cmdLookUp(hash string, logs bool) error {

	docker := client.NewDockerClient()

	// Check that RethinkDB is running
	if _, running, _ := container.Running(docker, "rethink"); !running {
		log.Error("RethinkDB is NOT running, starting now...")
		rethink, err := container.StartRethinkDB(docker, false)
		er.CheckError(err)
		rInfo, err := container.Inspect(docker, rethink.ID)
		er.CheckError(err)
		er.CheckError(database.TestConnection(rInfo.Node.Addr))
	}

	// Setup rethinkDB
	database.InitRethinkDB()

	if plugins.InstalledPluginsCheck(docker) {
		log.Debug("All enabled plugins are installed.")
	} else {
		// Prompt user to install all plugins?
		fmt.Println("All enabled plugins not installed would you like to install them now? (yes/no)")
		fmt.Println("[Warning] This can take a while if it is the first time you have ran Malice.")
		if util.AskForConfirmation() {
			plugins.UpdateAllPlugins(docker)
		}
	}

	/////////////////////////////
	// Write hash to the Database
	resp := database.WriteHashToDatabase(hash)

	plugins.RunIntelPlugins(docker, hash, resp.GeneratedKeys[0], true)

	return nil
}

// APILookUp is an API wrapper for cmdLookUp
func APILookUp(hash string) error {
	return cmdLookUp(hash, false)
}
