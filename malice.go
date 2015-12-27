package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/go-malice/commands"
	"github.com/blacktop/go-malice/config"
	"github.com/blacktop/go-malice/version"
	"github.com/codegangsta/cli"

	// "github.com/gorilla/handlers"
	// "github.com/jordan-wright/gophish/controllers"
	// "github.com/jordan-wright/gophish/models"
)

func init() {
	if config.Conf.Malice.Environment == "production" {
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

// AppHelpTemplate custom app help template
var AppHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

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

// CommandHelpTemplate custom command help template
var CommandHelpTemplate = `Usage: malice {{.Name}}{{if .Flags}} [OPTIONS]{{end}} [arg...]
{{.Usage}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .Flags}}

Options:
   {{range .Flags}}
   {{.}}{{end}}{{ end }}
`

func setDebugOutputLevel() {
	// TODO: I'm not really a fan of this method and really would rather
	// use -v / --verbose TBQH
	for _, f := range os.Args {
		if f == "-D" || f == "--debug" || f == "-debug" {
			log.SetLevel(log.DebugLevel)
		}
	}

	debugEnv := os.Getenv("MALICE_DEBUG")
	if debugEnv != "" {
		showDebug, err := strconv.ParseBool(debugEnv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing boolean value from MALICE_DEBUG: %s\n", err)
			os.Exit(1)
		}
		if showDebug {
			log.SetLevel(log.DebugLevel)
		}
	}
}

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	app := cli.NewApp()

	app.Name = "malice"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"

	app.Commands = commands.Commands
	app.CommandNotFound = cmdNotFound
	app.Usage = "Open Source Malware Analysis Framework"
	app.Version = version.FullVersion()
	app.EnableBashCompletion = true

	log.Debug("Malice Version: ", app.Version)

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Enable debug mode",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func cmdNotFound(c *cli.Context, command string) {
	log.Fatalf(
		"%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name,
		command,
		c.App.Name,
		os.Args[0],
	)
}
