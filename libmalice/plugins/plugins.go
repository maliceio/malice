package plugins

import (
	"bytes"
	"fmt"
	"log"

	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/maliceio/malice/config"
	"github.com/parnurzeal/gorequest"
)

// "github.com/pelletier/go-toml"

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
func InstallPlugin(plugin *config.Plugin) (err error) {

	var newPlugin = config.PluginConfiguration{
		[]config.Plugin{
			config.Plugin{
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
	assert(toml.NewEncoder(buf).Encode(newPlugin))
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
	plugins := config.Plug.Plugins
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

// filterPluginsByEnabled returns a map[string]plugin of plugins
// that work on the given mime type
func filterPluginsByMime(mime string) []config.Plugin {
	mimeMatch := []config.Plugin{}

	for _, plugin := range config.Plug.Plugins {
		if strings.Contains(plugin.Mime, mime) {
			mimeMatch = append(mimeMatch, plugin)
		}
	}
	return mimeMatch
}

// filterPluginsByEnabled returns a map[string]plugin of enalbed plugins
func filterPluginsByEnabled() []config.Plugin {
	enabled := []config.Plugin{}

	for _, plugin := range config.Plug.Plugins {
		if plugin.Enabled {
			enabled = append(enabled, plugin)
		}
	}
	return enabled
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
