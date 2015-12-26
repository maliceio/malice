package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/go-malice/config"
	"github.com/blacktop/go-malice/docker"
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
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "malice"
	app.Usage = "Open Source Malware Analysis Framework"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = "0.1.0"
	// app.Action = func(c *cli.Context) {
	// 	// fmt.Println(c.Bool("verbose"))
	// 	// fmt.Println(c.String("test"))
	// 	// println("Hello friend!")
	// }
	app.Commands = []cli.Command{
		{
			Name: "web",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Show more output",
				},
			},
			Usage: "start web api",
			Action: func(c *cli.Context) {
				// NOTE: Starting ELK container. I might use this to store data from Malice
				cont, err := docker.StartELK()
				if err != nil {
					fmt.Printf("StartELK error = %s\n", err)
				}
				log.WithFields(log.Fields{
					// "id":   cont.ID,
					"url":      "http://" + docker.GetIP(),
					"username": "admin",
					"password": "admin",
					"name":     cont.Name,
					"env":      config.Conf.Malice.Environment,
				}).Info("ELK Container Started")
				// searchFloom(c.Args(), c.Bool("verbose"))
				// Setup the global variables and settings
				// err := models.Setup()
				// if err != nil {
				// 	fmt.Println(err)
				// }
				// Start the web servers
				log.WithFields(log.Fields{
					"env": config.Conf.Malice.Environment,
					"url": "http://" + config.Conf.Malice.AdminURL,
				}).Info("Admin server started...")
				// go http.ListenAndServe(config.Conf.AdminURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAdminRouter()))
				log.WithFields(log.Fields{
					"env": config.Conf.Malice.Environment,
					"url": "http://" + config.Conf.Malice.URL,
				}).Info("Malice server started...")
				// http.ListenAndServe(config.Conf.PhishURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreatePhishingRouter()))

			},
			// BashComplete: func(c *cli.Context) {
			// 	// This will complete if no args are passed
			// 	if len(c.Args()) > 0 {
			// 		return
			// 	}
			// 	for _, t := range tasks {
			// 		fmt.Println(t)
			// 	}
			// },
		},
	}

	app.Run(os.Args)
}
