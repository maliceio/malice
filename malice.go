package main

import (
	"os"

	log "github.com/blacktop/docker-nsrl/nsrl/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/blacktop/docker-nsrl/nsrl/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/blacktop/go-malice/config"
	doc "github.com/blacktop/go-malice/docker"
	// "github.com/gorilla/handlers"
	// "github.com/jordan-wright/gophish/controllers"
	// "github.com/jordan-wright/gophish/models"
)

func init() {

	if config.Conf.Malice.Environment == "production" {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// log.SetFormatter(&logstash.LogstashFormatter{Type: "malice"})
	} else {
		// Log as ASCII formatter.
		log.SetFormatter(&log.TextFormatter{})
	}

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
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
				cont, err := doc.startELK()
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
