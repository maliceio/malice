package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
)

func init() {
	if config.Conf.Environment == "production" {
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

func cmdWebStart() error {

	// Setup the global variables and settings
	// err := models.Setup()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// Start the web servers
	log.WithFields(log.Fields{
		"env": config.Conf.Environment,
		"url": "http://" + config.Conf.Web.AdminURL,
	}).Info("Admin server started...")
	// go http.ListenAndServe(config.Config.AdminURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAdminRouter()))
	log.WithFields(log.Fields{
		"env": config.Conf.Environment,
		"url": "http://" + config.Conf.Web.URL,
	}).Info("Malice server started...")
	// http.ListenAndServe(config.Config.PhishURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreatePhishingRouter()))

	return nil
}

func cmdWebStop() error {

	// Setup the global variables and settings
	// err := models.Setup()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// Start the web servers
	log.WithFields(log.Fields{
		"env": config.Conf.Environment,
	}).Info("Admin server stopped...")
	// go http.ListenAndServe(config.Config.AdminURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAdminRouter()))
	log.WithFields(log.Fields{
		"env": config.Conf.Environment,
	}).Info("Malice server stopped...")
	// http.ListenAndServe(config.Config.PhishURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreatePhishingRouter()))

	return nil
}
