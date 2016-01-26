package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"os"
	"strings"
	// "github.com/pelletier/go-toml"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/docker/engine-api/types"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/data"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/parnurzeal/gorequest"
)

// Plugin represents a single plugin setting.
type Plugin struct {
	Name        string `toml:"name"`
	Enabled     bool   `toml:"enabled"`
	Category    string `toml:"category"`
	Description string `toml:"description"`
	Image       string `toml:"image"`
	Mime        string `toml:"mime"`
}

// PluginConfiguration represents the malice runtime plugins.
type PluginConfiguration struct {
	Plugins []Plugin `toml:"plugin"`
}

// Plug represents the Malice runtime configuration
var Plug PluginConfiguration

// Load plugins.toml into Plug var
func Load() {
	// Get the plugin file
	pluginPath := "./data/plugins.toml"
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		// er.CheckErrorNoStackWithMessage(err, "NOT FOUND")
		pluginPath = path.Join(maldirs.GetBaseDir(), "./plugins.toml")
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			pluginData, err := data.Asset("data/plugins.toml")
			er.CheckError(err)
			er.CheckError(ioutil.WriteFile(pluginPath, pluginData, 0644))
		}
	}
	log.Debug("Plugin Config: ", pluginPath)
	_, err := toml.DecodeFile(pluginPath, &Plug)
	er.CheckError(err)
}

// StartPlugin starts plugin
func (plugin Plugin) StartPlugin(client maldocker.Docker, sample string, logs bool) (types.ContainerJSONBase, error) {
	contJSON, err := client.StartContainer(sample, plugin.Name, plugin.Image, logs)
	er.CheckError(err)

	return contJSON, err
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

// UpdateAllPlugins performs a docker pull on all registered plugins checking for updates
func UpdateAllPlugins(client maldocker.Docker) {
	plugins := Plug.Plugins
	for _, plugin := range plugins {
		client.PullImage(plugin.Image, "latest")
	}
}

// UpdatePlugin performs a docker pull on all registered plugins checking for updates
func (plugin Plugin) UpdatePlugin(client maldocker.Docker) {
	client.PullImage(plugin.Image, "latest")
}

//InstallPlugin installs a new malice plugin
func InstallPlugin(plugin *Plugin) (err error) {

	var newPlugin = PluginConfiguration{
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

// ListEnabledPlugins lists enabled plugins
func ListEnabledPlugins(detail bool) {
	// TODO: Create a template for this kind of output : http://stackoverflow.com/questions/10747054/special-case-treatment-for-the-last-element-of-a-range-in-google-gos-text-templ
	enabled := filterPluginsByEnabled()
	if detail {
		ToMarkDownTable(enabled)
	} else {
		for _, plugin := range enabled {
			fmt.Println(plugin.Name)
		}
	}
}

// ListAllPlugins lists all plugins
func ListAllPlugins(detail bool) {
	plugins := Plug.Plugins
	if detail {
		ToMarkDownTable(plugins)
	} else {
		for _, plugin := range plugins {
			fmt.Println(plugin.Name)
		}
	}
}

// ToMarkDownTable prints plugins out as Markdown table
func ToMarkDownTable(plugins []Plugin) {
	table := clitable.New([]string{"Name", "Description", "Enabled", "Image", "Category", "Mime"})
	for _, plugin := range plugins {
		table.AddRow(map[string]interface{}{
			"Name":        plugin.Name,
			"Description": plugin.Description,
			"Enabled":     plugin.Enabled,
			"Image":       plugin.Image,
			"Category":    plugin.Category,
			"Mime":        plugin.Mime,
		})
	}
	table.Markdown = true
	table.Print()
}

// GetPluginByName will return plugin for the given name
func GetPluginByName(name string) Plugin {

	for _, plugin := range Plug.Plugins {
		if strings.EqualFold(plugin.Name, name) {
			return plugin
		}
	}

	return Plugin{}
}

// GetPluginsForMime will return all plugins that can consume the mime type file
func GetPluginsForMime(mime string) []Plugin {
	return filterPluginsByMime(mime)
}

// filterPluginsByEnabled returns a map[string]plugin of plugins
// that work on the given mime type
func filterPluginsByMime(mime string) []Plugin {
	mimeMatch := []Plugin{}

	for _, plugin := range filterPluginsByEnabled() {
		if strings.Contains(plugin.Mime, mime) || strings.Contains(plugin.Mime, "*") {
			mimeMatch = append(mimeMatch, plugin)
		}
	}
	return mimeMatch
}

// filterPluginsByEnabled returns a map[string]plugin of enalbed plugins
func filterPluginsByEnabled() []Plugin {
	enabled := []Plugin{}

	for _, plugin := range Plug.Plugins {
		if plugin.Enabled {
			enabled = append(enabled, plugin)
		}
	}
	return enabled
}
