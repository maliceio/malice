package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/plugins"
)

func init() {
	if config.Conf.Environment.Run == "production" {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.InfoLevel)
		// log.SetFormatter(&logstash.LogstashFormatter{Type: "malice"})
	} else {
		// Log as ASCII formatter.
		log.SetFormatter(&log.TextFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)
}

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
