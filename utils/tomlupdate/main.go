package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/maliceio/malice/config"
	"github.com/urfave/cli"
)

var (
	version   string
	buildtime string
)

func main() {
	app := cli.NewApp()

	app.Name = "tomlupdate"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = version + ", BuildTime: " + buildtime
	app.Compiled, _ = time.Parse("20060102", buildtime)
	app.Usage = "Update Version in Malice Config TOML"
	app.ArgsUsage = "malice VERSION"
	app.Action = func(c *cli.Context) error {
		if c.Args().Present() {
			// read VERSION
			ver := c.Args().First()
			// load config from TOML
			config.Load()
			// update config version from VERSION file
			config.Conf.Version = strings.TrimSpace(string(ver))

			buf := new(bytes.Buffer)
			if err := toml.NewEncoder(buf).Encode(config.Conf); err != nil {
				log.Fatal(err)
			}
			fmt.Println(buf.String())
		} else {
			log.Fatal(fmt.Errorf("Please supply a malice VERSION"))
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
