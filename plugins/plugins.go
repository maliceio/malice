package plugins

import (
	"bytes"
	"fmt"

	"os"
	"strings"
	// "github.com/pelletier/go-toml"
	"github.com/BurntSushi/toml"
	"github.com/fsouza/go-dockerclient"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/libmalice/errors"
	"github.com/maliceio/malice/libmalice/maldocker"
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

func init() {
	// Get the config file
	_, err := toml.DecodeFile("./plugins.toml", &Plug)
	er.CheckError(err)
	// fmt.Println(Plug)
}

// StartPlugin starts plugin
func (plugin Plugin) StartPlugin(sample string, logs bool) (cont *docker.Container, err error) {
	cont, err = maldocker.StartContainer(sample, plugin.Name, plugin.Image, logs)
	er.CheckError(err)

	// fmt.Println(cont.Name)
	return
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
		for idx, plugin := range enabled {
			fmt.Println("Name: ", plugin.Name)
			fmt.Println("Description: ", plugin.Description)
			fmt.Println("Enabled: ", plugin.Enabled)
			fmt.Println("Image: ", plugin.Image)
			fmt.Println("Category: ", plugin.Category)
			fmt.Println("Mime: ", plugin.Mime)
			if idx+1 < len(enabled) && len(enabled) != 1 {
				fmt.Println("---------------------")
			}
		}
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
		for idx, plugin := range plugins {
			fmt.Println("Name: ", plugin.Name)
			fmt.Println("Description: ", plugin.Description)
			fmt.Println("Enabled: ", plugin.Enabled)
			fmt.Println("Image: ", plugin.Image)
			fmt.Println("Category: ", plugin.Category)
			fmt.Println("Mime: ", plugin.Mime)
			if idx+1 < len(plugins) && len(plugins) != 1 {
				fmt.Println("---------------------")
			}
		}
	} else {
		for _, plugin := range plugins {
			fmt.Println(plugin.Name)
		}
	}
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
