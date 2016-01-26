package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/maliceio/malice/plugins"
)

func cmdListPlugins(all bool, detail bool) {
	if all {
		plugins.ListAllPlugins(detail)
	} else {
		plugins.ListEnabledPlugins(detail)
	}
	// enabled := plugins.GetEnabledPlugins()
	// fmt.Println(enabled)
}

func cmdInstallPlugin(name string) {

	testPlugin := plugins.Plugin{
		Name:        name,
		Enabled:     true,
		Category:    "test",
		Description: "This is a test plugin",
		Image:       "blacktop/test",
		Mime:        "image/png",
	}
	plugins.InstallPlugin(&testPlugin)
}

func cmdRemovePlugin() {

}

func cmdUpdatePlugin(name string, all bool) {
	docker := maldocker.NewDockerClient()
	if all {
		plugins.UpdateAllPlugins(docker)
	} else {
		if name == "" {
			log.Error("Please enter a valid plugin name.")
			os.Exit(1)
		}
		plugins.GetPluginByName(name).UpdatePlugin(docker)
	}
}
