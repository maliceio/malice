package commands

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var tasks = []string{"start", "stop"}

// Commands are the codegangsta/cli commands for Malice
var Commands = []cli.Command{
	{
		Name:        "scan",
		Usage:       "Scan a file",
		Description: "File to be scanned.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "logs",
				Usage: "Display the Logs of the Plugin containers",
			},
		},
		Action: func(c *cli.Context) error { return cmdScan(c.Args().First(), c.Bool("logs")) },
	},
	{
		Name:        "watch",
		Usage:       "Watch a folder",
		Description: "Folder to be watched.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "logs",
				Usage: "Display the Logs of the Plugin containers",
			},
		},
		Action: func(c *cli.Context) error { return cmdWatch(c.Args().First(), c.Bool("logs")) },
	},
	{
		Name:        "lookup",
		Usage:       "Look up a file hash",
		Description: "Hash to be queried.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "logs",
				Usage: "Display the Logs of the Plugin containers",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Present() {
				return cmdLookUp(c.Args().First(), c.Bool("logs"))
			}
			log.Fatal(fmt.Errorf("Please supply a MD5/SHA1 hash to query."))

			return nil
		},
	},
	{
		Name:        "elk",
		Usage:       "Start an ELK docker container",
		Description: "Argument is what port to bind to.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "logs",
				Usage: "Display the Logs from the ELK Container",
			},
		},
		Action: func(c *cli.Context) error { return cmdELK(c.Bool("logs")) },
	},
	{
		Name:  "web",
		Usage: "Start, Stop Web services",
		Subcommands: []cli.Command{
			{
				Name:   "start",
				Usage:  "start web application",
				Action: func(c *cli.Context) error { return cmdWebStart() },
			},
			{
				Name:   "stop",
				Usage:  "stop web application",
				Action: func(c *cli.Context) error { return cmdWebStop() },
			},
		},
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if len(c.Args()) > 0 {
				return
			}
			for _, t := range tasks {
				fmt.Println(t)
			}
		},
	},
	{
		Name:  "plugin",
		Usage: "List, Install or Remove Plugins",
		Subcommands: []cli.Command{
			{
				Name:  "list",
				Usage: "list enabled installed plugins",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "all",
						Usage: "display all installed plugins",
					},
					cli.BoolFlag{
						Name:  "detail,d",
						Usage: "display plugin details",
					},
				},
				Action: func(c *cli.Context) error { return cmdListPlugins(c.Bool("all"), c.Bool("detail")) },
			},
			{
				Name:   "install",
				Usage:  "install plugin",
				Action: func(c *cli.Context) error { return cmdInstallPlugin(c.Args().First()) },
			},
			{
				Name:   "remove",
				Usage:  "remove plugin",
				Action: func(c *cli.Context) error { return cmdRemovePlugin() },
			},
			{
				Name:  "update",
				Usage: "update plugin",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "all",
						Usage: "update all installed plugins",
					},
					cli.BoolFlag{
						Name:  "s,source",
						Usage: "update plugin from source repo",
					},
				},
				Action: func(c *cli.Context) error { return cmdUpdatePlugin(c.Args().First(), c.Bool("all"), c.Bool("source")) },
			},
		},
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if len(c.Args()) > 0 {
				return
			}
			for _, t := range tasks {
				fmt.Println(t)
			}
		},
	},
}

// CmdNotFound outputs a formatted command not found message
func CmdNotFound(c *cli.Context, command string) {
	log.Fatalf("%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, os.Args[0])
}
