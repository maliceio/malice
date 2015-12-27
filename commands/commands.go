package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var tasks = []string{"start", "stop"}

// Commands are the codegangsta/cli commands for Malice
var Commands = []cli.Command{
	{
		Name:        "elk",
		Usage:       "Start an ELK docker container",
		Description: "Argument is what port to bind to.",
		Action:      func(c *cli.Context) { cmdELK() },
	},
	{
		Name:    "web",
		Aliases: []string{"r"},
		Usage:   "options for task templates",
		Subcommands: []cli.Command{
			{
				Name:   "start",
				Usage:  "start web application",
				Action: func(c *cli.Context) { cmdWebStart() },
			},
			{
				Name:   "stop",
				Usage:  "stop web application",
				Action: func(c *cli.Context) { cmdWebStop() },
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
