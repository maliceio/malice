package maldocker

import (
	"os"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/engine-api/types"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"

	"regexp"

	log "github.com/Sirupsen/logrus"
)

// PullImage pulls docker image:tag
func (client *Docker) PullImage(id string, tag string) {

	pullOptions := types.ImagePullOptions{
		ImageID: id,
		Tag:     tag,
		// RegistryAuth: encodedAuth,
	}

	responseBody, err := client.Client.ImagePull(pullOptions, nil)
	defer responseBody.Close()
	er.CheckError(err)

	jsonmessage.DisplayJSONMessagesStream(responseBody, os.Stdout, os.Stdout.Fd(), true, nil)
}

// ImageExists returns APIImages images list and true
// if the image name exists, otherwise false.
func (client *Docker) ImageExists(name string) (types.Image, bool, error) {
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Searching for image: ", name)
	images, err := client.listImages(false)
	if err != nil {
		return types.Image{}, false, err
	}

	r := regexp.MustCompile(name)
	if len(images) != 0 {
		for _, image := range images {
			for _, tag := range image.RepoTags {
				if r.MatchString(tag) {
					log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Image FOUND: ", name)
					return image, true, nil
				}
			}
		}
	}

	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Image NOT Found: ", name)
	return types.Image{}, false, nil
}

func (client *Docker) listImages(all bool) ([]types.Image, error) {

	options := types.ImageListOptions{
		All: all,
		// MatchName string
		// Filters   filters.Args
	}
	imageList, err := client.Client.ImageList(options)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return imageList, nil
}
