package plugins

import (
	"bytes"
	"fmt"
	"sync"

	"os"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	"github.com/docker/engine-api/types/strslice"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/malice/docker/client/image"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/parnurzeal/gorequest"
)

// StartPlugin starts plugin
func (plugin Plugin) StartPlugin(docker *client.Docker, arg string, scanID string, logs bool, wg *sync.WaitGroup) {

	defer wg.Done()

	cmd := plugin.buildCmd(arg, logs)
	binds := []string{config.Conf.Docker.Binds}
	env := plugin.getPluginEnv()

	env = append(env, "MALICE_SCANID="+scanID)

	contJSON, err := container.Start(
		// err := container.Run(
		docker,
		cmd, //strslice.StrSlice{"-t", plugin.Cmd, arg},
		plugin.Name,
		plugin.Image,
		logs,
		binds,
		nil,
		[]string{"rethink"},
		env,
	)
	log.WithFields(log.Fields{
		"name": contJSON.Name,
		"env":  config.Conf.Environment.Run,
	}).Debug("Plugin Container Started")

	defer func() {
		er.CheckError(container.Remove(docker, contJSON.ID, true, false, true))
		log.WithFields(log.Fields{
			"name": contJSON.Name,
			"env":  config.Conf.Environment.Run,
		}).Debug("Plugin Container Removed")
	}()

	er.CheckError(err)
}

func (plugin Plugin) buildCmd(args string, logs bool) strslice.StrSlice {

	cmdStr := strslice.StrSlice{}
	if logs {
		cmdStr = append(cmdStr, "-t")
	}
	if plugin.Cmd != "" {
		cmdStr = append(cmdStr, plugin.Cmd)
	}
	cmdStr = append(cmdStr, args)

	return cmdStr
}

// RunIntelPlugins run all Intel plugins
func RunIntelPlugins(docker *client.Docker, hash string, scanID string, logs bool) {

	intelPlugins := GetIntelPlugins(true)

	var wg sync.WaitGroup
	wg.Add(len(intelPlugins))

	for _, plugin := range intelPlugins {
		go plugin.StartPlugin(docker, hash, scanID, logs, &wg)
	}
	wg.Wait()
}

func (plugin *Plugin) getPluginEnv() []string {
	var env []string
	for _, pluginEnv := range plugin.Env {
		env = append(env, fmt.Sprintf("%s=%s", pluginEnv, os.Getenv(pluginEnv)))
	}
	return env
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
func InstalledPluginsCheck(docker *client.Docker) bool {
	for _, plugin := range getEnabled(Plugs.Plugins) {
		if _, exists, _ := image.Exists(docker, plugin.Image); !exists {
			return false
		}
	}
	return true
}

// UpdatePlugin performs a docker pull on all registered plugins checking for updates
func (plugin Plugin) UpdatePlugin(docker *client.Docker) {
	image.Pull(docker, plugin.Image, "latest")
}

// UpdatePluginFromRepository performs a docker build on a plugins remote repository
func (plugin Plugin) UpdatePluginFromRepository(docker *client.Docker) {

	log.Info("[Building Plugin from Source] ===> ", plugin.Name)

	var buildArgs map[string]string
	var quiet = false

	tags := []string{"malice/" + plugin.Name + ":latest"}

	if config.Conf.Proxy.Enable {
		buildArgs = runconfigopts.ConvertKVStringsToMap([]string{
			"HTTP_PROXY=" + config.Conf.Proxy.HTTP,
			"HTTPS_PROXY=" + config.Conf.Proxy.HTTPS,
		})
	} else {
		buildArgs = nil
	}

	labels := runconfigopts.ConvertKVStringsToMap([]string{"io.malice.plugin.installed.from=repository"})

	image.Build(docker, plugin.Repository, tags, buildArgs, labels, quiet)
}

// UpdateAllPlugins performs a docker pull on all registered plugins checking for updates
func UpdateAllPlugins(docker *client.Docker) {
	plugins := Plugs.Plugins
	for _, plugin := range plugins {
		fmt.Println("[Updating Plugin] ===> ", plugin.Name)
		if plugin.Build {
			plugin.UpdatePluginFromRepository(docker)
		} else {
			image.Pull(docker, plugin.Image, "latest")
		}
	}
}

// UpdateAllPluginsFromSource performs a docker build on a plugins remote repository on all registered plugins
func UpdateAllPluginsFromSource(docker *client.Docker) {
	plugins := Plugs.Plugins
	for _, plugin := range plugins {
		fmt.Println("[Updating Plugin from Source] ===> ", plugin.Name)
		plugin.UpdatePluginFromRepository(docker)
	}
}
