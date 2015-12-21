package main

import (
	// "fmt"
	// "gophish-master/models"
	"log"
	"os"

	"github.com/blacktop/go-malice/config"
	// "github.com/gorilla/handlers"
	// "github.com/jordan-wright/gophish/controllers"
	// "github.com/jordan-wright/gophish/models"
)

// Logger - is logger
var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// Setup the global variables and settings
	// err := models.Setup()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// Start the web servers
	Logger.Printf("Admin server started at http://%s\n", config.Conf.Malice.AdminURL)
	// go http.ListenAndServe(config.Conf.AdminURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateAdminRouter()))
	Logger.Printf("Malice server started at http://%s\n", config.Conf.Malice.URL)
	// http.ListenAndServe(config.Conf.PhishURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreatePhishingRouter()))
	println(config.Conf.Malice.Email.Host)
	println(config.Conf.Malice.Email.Port)
	println(config.Conf.Malice.DB.Path)
}
