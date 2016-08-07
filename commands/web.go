package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
)

func cmdWebStart() error {

	// Setup the global variables and settings
	// err := models.Setup()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// Start the web servers
	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
		"url": "http://" + config.Conf.Web.AdminURL,
	}).Info("Admin server started...")
	// go http.ListenAndServe(config.Config.AdminURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAdminRouter()))
	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
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
		"env": config.Conf.Environment.Run,
	}).Info("Admin server stopped...")
	// go http.ListenAndServe(config.Config.AdminURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAdminRouter()))
	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Info("Malice server stopped...")
	// http.ListenAndServe(config.Config.PhishURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreatePhishingRouter()))

	return nil
}
