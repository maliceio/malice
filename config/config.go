package config

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
	"log"
)

// Config represents the configuration information.
type Config struct {
	Malice struct {
		URL      string     `yaml:"url"`
		AdminURL string     `yaml:"admin_url"`
		Email    SMTPServer `yaml:"email"`
		DB       Database   `yaml:"db"`
	}
}

// SMTPServer represents the Email configuration details
type SMTPServer struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Database represents the Database configuration details
type Database struct {
	Path string `yaml:"path"`
}

// type Config struct {
// 	AdminURL string     `json:"admin_url"`
// 	PhishURL string     `json:"phish_url"`
// 	SMTP     SMTPServer `json:"smtp"`
// 	DBPath   string     `json:"dbpath"`
// }

var Conf Config

func init() {
	// Get the config file
	config_file, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	yaml.Unmarshal(config_file, &Conf)
}
