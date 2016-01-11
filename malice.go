package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/maliceio/malice/commands"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/version"
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
		log.SetLevel(log.InfoLevel)
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
	app.CommandNotFound = commands.CmdNotFound
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

// NEXT: Get plugins to return JSON or Markdown tables
// NEXT: Get plugins to POST JSON to malice webhook which pipes it to ELK Container
// NEXT: Binpack config/plugins.toml into Malice to write out to .malice dir on first run
// NEXT: Check if file already exists then display stored results
// NEXT: Rewrite to reuse containers instead of spawning new one all the time (will speed up AV)
// NEXT: Rewrite Plugins into goroutines for speed
// NEXT: Cleanup code, docker/plugin/container/ etc utils !!!
