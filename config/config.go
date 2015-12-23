package config

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
	"log"
)

// Config represents the configuration information.
type Config struct {
	Malice struct {
		URL         string
		AdminURL    string     `yaml:"admin_url"`
		Email       SMTPServer `yaml:"email"`
		DB          Database   `yaml:"db"`
		Environment string     `yaml:"env"`
		Docker      Docker
	}
}

// SMTPServer represents the Email configuration details
type SMTPServer struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Database represents the Database configuration details
type Database struct {
	Path string
}

// Docker represents the Docker configuration details
type Docker struct {
	Endpoint string
}

// Conf represents the Malice runtime configuration
var Conf Config

func init() {
	// Get the config file
	config, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	yaml.Unmarshal(config, &Conf)
}
