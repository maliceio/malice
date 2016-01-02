package maldocker

import (
	"github.com/maliceio/malice/config"

	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

func init() {
	if config.Conf.Environment.Run == "production" {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.InfoLevel)
		// log.SetFormatter(&logstash.LogstashFormatter{Type: "malice"})
	} else {
		// Log as ASCII formatter.
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)
}

// PullImage pulls docker image:tag
func PullImage(client *docker.Client, name string, tag string) (err error) {

	opts := docker.PullImageOptions{
		Repository: name,
		// Registry      string
		Tag:          tag,
		OutputStream: os.Stdout,
		// RawJSONStream: true,
	}

	auth := docker.AuthConfiguration{
	// Username      string `json:"username,omitempty"`
	// Password      string `json:"password,omitempty"`
	// Email         string `json:"email,omitempty"`
	// ServerAddress string `json:"serveraddress,omitempty"`
	}

	err = client.PullImage(opts, auth)

	return
}

// ImageExists returns APIImages images list and true
// if the image name exists, otherwise false.
func ImageExists(client *docker.Client, name string) (*docker.APIImages, bool, error) {
	return ParseImages(client, name)
}

// ParseImages parses the images
func ParseImages(client *docker.Client, name string) (*docker.APIImages, bool, error) {
	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Debug("Searching for image: ", name)
	images, err := listImages(client)
	if err != nil {
		return nil, false, err
	}
	r := regexp.MustCompile(name)

	if len(images) != 0 {
		for _, image := range images {
			for _, tag := range image.RepoTags {
				if r.MatchString(tag) {
					log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Image FOUND: ", name)
					return &image, true, nil
				}
			}
		}
	}

	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Debug("Image NOT Found: ", name)

	return nil, false, nil
	// return docker.APIImages{}, false, nil
}

func listImages(client *docker.Client) ([]docker.APIImages, error) {
	var images []docker.APIImages

	imageList, err := client.ListImages(docker.ListImagesOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, image := range imageList {
		images = append(images, image)
	}

	return images, nil
}
