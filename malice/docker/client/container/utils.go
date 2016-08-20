package container

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	er "github.com/maliceio/malice/malice/errors"

	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/image"
	"github.com/maliceio/malice/malice/docker/client/network"
	"github.com/maliceio/malice/malice/docker/client/volume"
)

func checkContainerRequirements(docker *client.Docker, containerName, img string) {
	// Check for existance of malice network
	if _, exists, _ := network.Exists(docker, "malice"); !exists {
		log.WithFields(log.Fields{
			"network": "malice",
			"exisits": exists,
			"env":     config.Conf.Environment.Run,
		}).Error("Network malice does not exist, creating now...")
		_, err := network.Create(docker, "malice")
		er.CheckError(err)
	}
	// Check for existance of malice volume
	if _, exists, _ := volume.Exists(docker, "malice"); !exists {
		log.Debug("Volume malice not found.")
		er.CheckError(volume.Create(docker, "malice", "local", nil))
	}
	log.Debug("Volume malice found.")
	// Check that the container isn't already running
	if _, exists, _ := Exists(docker, containerName); exists {
		log.WithFields(log.Fields{
			"exisits": exists,
			"name":    containerName,
			"env":     config.Conf.Environment.Run,
		}).Error("Container is already running...")
		os.Exit(0)
	}
	// Check that we have already pulled the image
	if _, exists, _ := image.Exists(docker, img); exists {
		log.WithFields(log.Fields{
			"exisits": exists,
			"env":     config.Conf.Environment.Run,
		}).Debugf("Image `%s` already pulled.", img)
	} else {
		log.WithFields(log.Fields{
			"exisits": exists,
			"env":     config.Conf.Environment.Run}).Debugf("Pulling Image `%s`", img)
		image.Pull(docker, img, "latest")
	}
}

// ErrConnectionFailed is an error raised when the connection between the client and the server failed.
var ErrConnectionFailed = errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")

// ErrorConnectionFailed returns an error with host in the error message when connection to docker daemon failed.
func ErrorConnectionFailed(host string) error {
	return fmt.Errorf("Cannot connect to the Docker daemon at %s. Is the docker daemon running?", host)
}

// getExitCode performs an inspect on the container. It returns
// the running state and the exit code.
func getExitCode(docker *client.Docker, ctx context.Context, containerID string) (bool, int, error) {
	c, err := Inspect(docker, containerID)
	if err != nil {
		// If we can't connect, then the daemon probably died.
		if err != ErrConnectionFailed {
			return false, -1, err
		}
		return false, -1, nil
	}
	return c.State.Running, c.State.ExitCode, nil
}
