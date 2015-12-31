package config

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	// "github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

// Plugins represents the configuration information.
type Plugins struct {
	Bin BinaryPlugin   `yaml:"binary"`
	Doc DocumantPlugin `yaml:"document"`
}

// BinaryPlugin represents the Email configuration details
type BinaryPlugin struct {
	Name struct {
		Enabled string `yaml:"enabled"`
		Image   string `yaml:"image"`
	}
}

// DocumantPlugin represents the Database configuration details
type DocumantPlugin struct {
	Name struct {
		Enabled string `yaml:"enabled"`
		Image   string `yaml:"image"`
	}
}

type tomlConfig struct {
	Title   string
	Owner   ownerInfo
	DB      database `toml:"database"`
	Plugins map[string]plugin
	Clients clients
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
	Bio  string
	DOB  time.Time
}

type database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type plugin struct {
	Enabled     bool
	Category    string
	Description string
	Image       string
	Mime        string
}

type clients struct {
	Data  [][]interface{}
	Hosts []string
}

// Plugin represents the Malice regiestered Plugins
var Plugin Plugins

// TConf represents the Malice regiestered Plugins
var TConf tomlConfig

func init() {
	// Get the config file
	plugins, err := ioutil.ReadFile("./plugins.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	yaml.Unmarshal(plugins, &Plugin)

	// ###############################################################
	if _, err := toml.DecodeFile("./config.toml", &TConf); err != nil {
		log.Fatalf("error: %v", err)
	}
}
