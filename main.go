package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/maliceio/malice/commands"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/logger"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/maliceio/malice/plugins"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	logger.Init(version)
	setDebugOutputLevel()
	config.Load(version)
	plugins.Load()
}

func setDebugOutputLevel() {
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

// Init initializes Malice
func Init() {
	maldirs.MakeDirs()
}

func main() {
	log.Debugf("Using %d PROCS", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	Init()
	// setDebugOutputLevel()
	cli.AppHelpTemplate = commands.AppHelpTemplate
	cli.CommandHelpTemplate = commands.CommandHelpTemplate
	app := cli.NewApp()

	app.Name = "malice"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"

	app.Commands = commands.Commands
	app.CommandNotFound = commands.CmdNotFound
	app.Usage = "Open Source Malware Analysis Framework"
	app.Version = fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
	app.Copyright = "Copyright (c) 2013 - 2016 'blacktop' Joshua Maine"
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
		log.Fatalln(err)
	}
}
