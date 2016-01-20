package maldocker

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/jsonmessage"
	dockerAPI "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/maliceio/malice/config"

	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

// PullImage pulls docker image:tag
func PullImage(imageName string, imageTag string) (err error) {

	cli, err := dockerAPI.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(options)
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		fmt.Println(c.ID)
	}

	pullOptions := types.ImagePullOptions{
		ImageID: imageName,
		Tag:     imageTag,
		// RegistryAuth: encodedAuth,
	}

	responseBody, err := cli.ImagePull(pullOptions, nil)
	if err != nil {
		panic(err)
	}
	defer responseBody.Close()

	jsonmessage.DisplayJSONMessagesStream(responseBody, os.Stdout, os.Stdout.Fd(), true, nil)

	return err
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
