package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

var header = `# Copyright (c) 2013 - 2017 blacktop Joshua Maine, All Rights Reserved.
# See LICENSE for license information.

#######################################################################
# MALICE Configuration ################################################
#######################################################################

`

func main() {

	var path string
	app := cli.NewApp()

	app.Name = "tomlupdate"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = version + ", BuildTime: " + buildtime
	app.Compiled, _ = time.Parse("20060102", buildtime)
	app.Usage = "Update Version in Malice Config TOML"
	app.ArgsUsage = "malice VERSION"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "path",
			Value:       "",
			Usage:       "path to malice config TOML",
			Destination: &path,
		},
	}
	app.Action = func(c *cli.Context) error {
		if path == "" {
			log.Fatal(fmt.Errorf("please supply `--path` of config.toml flle"))
		}
		if c.Args().Present() {
			// read VERSION
			ver := strings.TrimSpace(string(c.Args().First()))
			// load config from TOML
			config.LoadFromToml("config/config.toml", ver)
			// update config version from VERSION file
			config.Conf.Version = ver

			buf := new(bytes.Buffer)
			if err := toml.NewEncoder(buf).Encode(config.Conf); err != nil {
				panic(err)
			}
			fmt.Println(buf.String())
			// open plugin config file
			configPath, err := filepath.Abs(path)
			if err != nil {
				log.Fatal(err)
			}
			f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			// write new plugin to installed plugin config
			if _, err = f.WriteString(header + buf.String()); err != nil {
				panic(err)
			}
		} else {
			log.Fatal(fmt.Errorf("Please supply a malice VERSION"))
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
