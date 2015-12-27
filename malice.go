package main

import (
	"os"

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

func main() {
	cli.AppHelpTemplate = commands.AppHelpTemplate
	cli.CommandHelpTemplate = commands.CommandHelpTemplate
	app := cli.NewApp()

	app.Name = "malice"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"

	app.Commands = commands.Commands
	app.CommandNotFound = cmdNotFound
	app.Usage = "Open Source Malware Analysis Framework"
	app.Version = version.FullVersion()
	// app.EnableBashCompletion = true

	log.Debug("Malice Version: ", app.Version)

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			EnvVar: "MALICE_DEBUG",
			Name:   "debug, D",
			Usage:  "Enable debug mode",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func cmdNotFound(c *cli.Context, command string) {
	log.Fatalf("%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, os.Args[0])
}
