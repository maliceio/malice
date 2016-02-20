package plugins

import (
	"bytes"
	"fmt"

	"os"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/strslice"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/parnurzeal/gorequest"
)

// StartPlugin starts plugin
func (plugin Plugin) StartPlugin(client *maldocker.Docker, sample string, logs bool) (types.ContainerJSONBase, error) {
	contJSON, err := client.StartContainer(strslice.New("-t", sample), plugin.Name, plugin.Image, logs)
	er.CheckError(err)

	return contJSON, err
}

// RunIntelPlugins run all Intel plugins
func RunIntelPlugins(client *maldocker.Docker, hash string, logs bool) {
	var cont types.ContainerJSONBase
	var err error
	for _, plugin := range GetIntelPlugins(true) {
		if plugin.Cmd != "" {
			cont, err = client.StartContainer(strslice.New("-t", plugin.Cmd, hash), plugin.Name, plugin.Image, logs)
			er.CheckError(err)
		} else {
			cont, err = client.StartContainer(strslice.New("-t", hash), plugin.Name, plugin.Image, logs)
			er.CheckError(err)
		}
		log.WithFields(log.Fields{
			"id": cont.ID,
			"ip": client.GetIP(),
			// "url":      "http://" + maldocker.GetIP(),
			"name": cont.Name,
			"env":  config.Conf.Environment.Run,
		}).Debug("Plugin Container Started")

		err := client.RemoveContainer(cont, false, false, false)
		er.CheckError(err)
	}
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

// PostResults post plugin results to Malice Webhook
func PostResults(url string, resultJSON []byte, taskID string) {
	request := gorequest.New()
	if config.Conf.Proxy.Enable {
		request = gorequest.New().Proxy(config.Conf.Proxy.HTTP)
	}
	request.Post(url).
		Set("Task", taskID).
		Send(resultJSON).
		End(printStatus)
}

//InstallPlugin installs a new malice plugin
func InstallPlugin(plugin *Plugin) (err error) {

	var newPlugin = Configuration{
		[]Plugin{
			Plugin{
				Name:        plugin.Name,
				Enabled:     plugin.Enabled,
				Category:    plugin.Category,
				Description: plugin.Description,
				Image:       plugin.Image,
				Mime:        plugin.Mime,
			},
		},
	}

	buf := new(bytes.Buffer)
	er.CheckError(toml.NewEncoder(buf).Encode(newPlugin))
	fmt.Println(buf.String())
	// open plugin config file
	f, err := os.OpenFile("./plugins.toml", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// write new plugin to installed plugin config
	if _, err = f.WriteString("\n" + buf.String()); err != nil {
		panic(err)
	}
	return
}

// InstalledPluginsCheck checks that all enabled plugins are installed
func InstalledPluginsCheck(client *maldocker.Docker) bool {
	for _, plugin := range getEnabled(Plugs.Plugins) {
		if _, exists, _ := client.ImageExists(plugin.Image); !exists {
			return false
		}
	}
	return true
}

// UpdatePlugin performs a docker pull on all registered plugins checking for updates
func (plugin Plugin) UpdatePlugin(client *maldocker.Docker) {
	client.PullImage(plugin.Image, "latest")
}

// UpdateAllPlugins performs a docker pull on all registered plugins checking for updates
func UpdateAllPlugins(client *maldocker.Docker) {
	plugins := Plugs.Plugins
	for _, plugin := range plugins {
		client.PullImage(plugin.Image, "latest")
	}
}
