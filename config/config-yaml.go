package config

//
// import (
// 	"gopkg.in/yaml.v2"
//
// 	"io/ioutil"
// 	"log"
// )
//
// // Config represents the configuration information.
// type Config struct {
// 	Malice struct {
// 		URL         string
// 		AdminURL    string     `yaml:"admin_url"`
// 		Email       SMTPServer `yaml:"email"`
// 		DB          Database   `yaml:"db"`
// 		Environment string     `yaml:"env"`
// 		Docker      Docker
// 		Log         Logger `yaml:"log"`
// 	}
// }
//
// // SMTPServer represents the Email configuration details
// type SMTPServer struct {
// 	Host     string
// 	Port     string
// 	User     string
// 	Password string
// }
//
// // Database represents the Database configuration details
// type Database struct {
// 	Path string
// }
//
// // Docker represents the Docker configuration details
// type Docker struct {
// 	Name     string `yaml:"machine-name"`
// 	Endpoint string
// }
//
// // Logger represents the Logger configuration details
// type Logger struct {
// 	// Filename is the file to write logs to.  Backup log files will be retained
// 	// in the same directory.  It uses <processname>-lumberjack.log in
// 	// os.TempDir() if empty.
// 	Filename string `json:"filename" yaml:"filename"`
//
// 	// MaxSize is the maximum size in megabytes of the log file before it gets
// 	// rotated. It defaults to 100 megabytes.
// 	MaxSize int `json:"maxsize" yaml:"maxsize"`
//
// 	// MaxAge is the maximum number of days to retain old log files based on the
// 	// timestamp encoded in their filename.  Note that a day is defined as 24
// 	// hours and may not exactly correspond to calendar days due to daylight
// 	// savings, leap seconds, etc. The default is not to remove old log files
// 	// based on age.
// 	MaxAge int `json:"maxage" yaml:"maxage"`
//
// 	// MaxBackups is the maximum number of old log files to retain.  The default
// 	// is to retain all old log files (though MaxAge may still cause them to get
// 	// deleted.)
// 	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`
//
// 	// LocalTime determines if the time used for formatting the timestamps in
// 	// backup files is the computer's local time.  The default is to use UTC
// 	// time.
// 	LocalTime bool `json:"localtime" yaml:"localtime"`
// 	// contains filtered or unexported fields
// }
//
// // Conf represents the Malice runtime configuration
// var Conf Config
//
// func init() {
// 	// Get the config file
// 	config, err := ioutil.ReadFile("./config.yml")
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}
// 	yaml.Unmarshal(config, &Conf)
// }
