package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/plugins"
)

func cmdListPlugins(all bool, detail bool) error {
	if all {
		plugins.ListAllPlugins(detail)
	} else {
		plugins.ListEnabledPlugins(detail)
	}

	// TODO: Add ability to list malice plugins not installed

	// docker := client.NewDockerClient()
	// err := docker.SearchImages("malice")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// enabled := plugins.GetEnabledPlugins()
	// fmt.Println(enabled)
	return nil
}

func cmdInstallPlugin(name string) error {

	testPlugin := plugins.Plugin{
		Name:        name,
		Enabled:     true,
		Category:    "test",
		Description: "This is a test plugin",
		Image:       "blacktop/test",
		Mime:        "image/png",
	}
	plugins.InstallPlugin(&testPlugin)

	return nil
}

func cmdRemovePlugin() error {
	return nil
}

func cmdUpdatePlugin(name string, all bool, source bool) error {
	docker := client.NewDockerClient()
	if all {
		plugins.UpdateAllPlugins(docker)
	} else {
		if name == "" {
			log.Error("Please enter a valid plugin name.")
			os.Exit(1)
		}
		if source {
			plugins.GetPluginByName(name).UpdatePluginFromRepository(docker)
		} else {
			plugins.GetPluginByName(name).UpdatePlugin(docker)
		}
	}
	return nil
}
