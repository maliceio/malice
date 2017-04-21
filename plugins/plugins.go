package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"

	"os"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/strslice"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/malice/docker/client/image"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/parnurzeal/gorequest"
)

// StartPlugin starts plugin
func (plugin Plugin) StartPlugin(docker *client.Docker, arg string, scanID string, logs bool, wg *sync.WaitGroup) {

	defer wg.Done()

	cmd := plugin.buildCmd(arg, logs)
	binds := []string{config.Conf.Docker.Binds} // []string{maldirs.GetSampledsDir() + ":/malware:ro"},
	env := plugin.getPluginEnv()

	env = append(env, "MALICE_SCANID="+scanID)
	log.WithFields(log.Fields{
		"name": plugin.Name,
		"env":  config.Conf.Environment.Run,
	}).Debug("env: ", env)
	// env = append(env, "MALICE_ELASTICSEARCH="+utils.Getopt("MALICE_ELASTICSEARCH", getDbAddr()))

	contJSON, err := container.Start(
		docker,       // docker *client.Docker,
		cmd,          // cmd strslice.StrSlice,
		plugin.Name,  // name string,
		plugin.Image, // image string,
		logs,         // logs bool,
		binds,        // binds []string,
		nil,          // portBindings nat.PortMap,
		[]string{config.Conf.Docker.Links}, // links []string,
		env, // env []string,
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

// getDbAddr gets address of DB server
func getDbAddr() string {
	return "localhost"
}

// buildCmd creates plugin run command
func (plugin Plugin) buildCmd(args string, logs bool) strslice.StrSlice {

	cmdStr := strslice.StrSlice{}
	if plugin.APIKey != "" {
		cmdStr = append(cmdStr, "--api", plugin.APIKey)
	}
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

	hashType, _ := utils.GetHashType(hash)

	log.Debug("Looking for Intel plugins...")
	intelPlugins := GetIntelPlugins(hashType, true)
	log.Debug("Found these plugins: ")
	for _, plugin := range intelPlugins {
		log.Debugf(" - %v", plugin.Name)
	}

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
		if os.Getenv(pluginEnv) != "" {
			env = append(env, fmt.Sprintf("%s=%s", pluginEnv, os.Getenv(pluginEnv)))
		}
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
	// fmt.Println(buf.String())

	// open plugin config file
	configPath := path.Join(maldirs.GetPluginsDir(), "./plugins.toml")
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0600)
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

// DeletePlugin deletes a plugin
func DeletePlugin(name string) error {

	for i, plugin := range Plugs.Plugins {
		if strings.EqualFold(plugin.Name, name) {
			Plugs.Plugins = append(Plugs.Plugins[:i], Plugs.Plugins[i+1:]...)
			break
		}
	}

	buf := new(bytes.Buffer)
	er.CheckError(toml.NewEncoder(buf).Encode(Plugs))

	// open plugin config file
	configPath := path.Join(maldirs.GetPluginsDir(), "./plugins.toml")
	err := ioutil.WriteFile(configPath, buf.Bytes(), 0644)
	return err
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

	var buildArgs map[string]*string
	var quiet = false

	tags := []string{"malice/" + plugin.Name + ":latest"}

	if config.Conf.Proxy.Enable {
		buildArgs = runconfigopts.ConvertKVStringsToMapWithNil([]string{
			"HTTP_PROXY=" + config.Conf.Proxy.HTTP,
			"HTTPS_PROXY=" + config.Conf.Proxy.HTTPS,
		})
	} else {
		buildArgs = nil
	}

	labels := runconfigopts.ConvertKVStringsToMap([]string{"io.malice.plugin.installed.from=repository"})

	image.Build(docker, plugin.Repository, tags, buildArgs, labels, quiet)
}

// UpdateEnabledPlugins performs a docker pull on all enabled plugins checking for updates
func UpdateEnabledPlugins(docker *client.Docker) {
	// Pull busybox (used to copy samples to malice volume)
	image.Pull(docker, "busybox", "latest")
	// Pull blacktop/elk (used to store malice scan results data)
	image.Pull(docker, config.Conf.DB.Image, "latest")
	// Pull all enabled malice plugin images
	for _, plugin := range GetEnabledPlugins() {
		fmt.Println("[Updating Plugin] ===> ", plugin.Name)
		if plugin.Build {
			plugin.UpdatePluginFromRepository(docker)
		} else {
			image.Pull(docker, plugin.Image, "latest")
		}
	}
}

// UpdateAllPlugins performs a docker pull on all registered plugins checking for updates
func UpdateAllPlugins(docker *client.Docker) {
	// Pull busybox (used to copy samples to malice volume)
	image.Pull(docker, "busybox", "latest")
	// Pull blacktop/elk (used to store malice scan results data)
	image.Pull(docker, config.Conf.DB.Image, "latest")
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
