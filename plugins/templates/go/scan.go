package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/maliceio/go-plugin-utils/database/elasticsearch"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/urfave/cli"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

const (
	name     = "{{ plugin_name }}"
	category = "{{ plugin_category }}"
)

type pluginResults struct {
	ID   string      `json:"id" structs:"id,omitempty"`
	Data resultsData `json:"{{ plugin_name }}" structs:"{{ plugin_name }}"`
}

type {{ plugin_name }} struct {
	Results resultsData `json:"{{ plugin_name }}"`
}

type resultsData struct {
	Data   []string         `json:"data" structs:"data"`
}

func printMarkDownTable(f {{ plugin_name }}) {
	fmt.Printf("#### {{ plugin_name }}\n\n")
}

func parsePluginOutput(pluginOutput string, all bool) resultsData {

	keepLines := []string{}
	results := resultsData{}

	lines := strings.Split(pluginOutput, "\n")

	// remove empty lines
	for _, line := range lines {
		if len(strings.TrimSpace(line)) != 0 {
			keepLines = append(keepLines, strings.TrimSpace(line))
		}
	}
	// build results data
	for i := 0; i < len(keepLines); i++ {

	}

	return results
}

// scanFile scans file
func scanFile(path string, all bool) pluginResults {
	pResults := pluginResults{}
	pResults.Results = parsePluginOutput(utils.RunCommand("/usr/bin/plugin_binary", "-flag", path))

	return pResults
}

var appHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

{{.Usage}}

Version: {{.Version}}{{if or .Author .Email}}

Author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

func main() {
	cli.AppHelpTemplate = appHelpTemplate
	app := cli.NewApp()
	app.Name = "{{ plugin_name }}"
	app.Author = "{{ author }}"
	app.Email = "{{ email }}"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "{{ usage }}"
	var table bool
	var all bool
	var elasticsearch string
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "",
			Usage:       "elasticsearch address for Malice to store results",
			EnvVar:      "MALICE_ELASTICSEARCH",
			Destination: &elasticsearch,
		},
		cli.BoolFlag{
			Name:   "post, p",
			Usage:  "POST results to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
		cli.BoolFlag{
			Name:        "table, t",
			Usage:       "output as Markdown table",
			Destination: &table,
		},
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "output ascii/utf-16 strings",
			Destination: &all,
		},
	}
	app.ArgsUsage = "FILE to scan with {{ plugin_name }}"
	app.Action = func(c *cli.Context) error {
		if c.Args().Present() {
			path := c.Args().First()
			// Check that file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				utils.Assert(err)
			}

			if c.Bool("verbose") {
				log.SetLevel(log.DebugLevel)
			}

			pluginOut := scanFile(path, all)

			// upsert into Database
			elasticsearch.InitElasticSearch()
			elasticsearch.WritePluginResultsToDatabase(elasticsearch.PluginResults{
				ID:       utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
				Name:     name,
				Category: category,
				Data:     structs.Map(pluginOut.Results),
			})

			if table {
				printMarkDownTable(pluginOut)
			} else {
				pluginJSON, err := json.Marshal(pluginOut)
				utils.Assert(err)
				fmt.Println(string(pluginJSON))
			}
		} else {
			log.Fatal(fmt.Errorf("Please supply a file to scan with {{ plugin_name }}"))
		}
		return nil
	}

	err := app.Run(os.Args)
	utils.Assert(err)
}
