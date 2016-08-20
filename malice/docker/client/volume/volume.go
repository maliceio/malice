package volume

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	"golang.org/x/net/context"
)

// Exists returns type.Volume and true
// if the volume name exists, otherwise false.
func Exists(docker *client.Docker, name string) (*types.Volume, bool, error) {
	return parseVolumes(docker, name, true)
}

// Create creates a docker volume with the given name
// returns: error
func Create(docker *client.Docker, name, driver string, labels []string) error {
	volReq := types.VolumeCreateRequest{
		Driver: driver,
		// DriverOpts: opts.driverOpts.GetAll(),
		Name:   name,
		Labels: runconfigopts.ConvertKVStringsToMap(labels),
	}

	vol, err := docker.Client.VolumeCreate(context.Background(), volReq)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"env": config.Conf.Environment.Run,
	}).Info("Created Volume: ", vol.Name)

	return nil
}

// ParseVolumes parses the volumes
func parseVolumes(docker *client.Docker, name string, all bool) (*types.Volume, bool, error) {
	// list volumes
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Searching for volume: ", name)
	volumes, err := List(docker, all)
	if err != nil {
		return nil, false, err
	}
	// locate docker volume that matches name
	r := regexp.MustCompile(name)
	if len(volumes.Volumes) != 0 {
		for _, volume := range volumes.Volumes {
			if r.MatchString(volume.Name) {
				log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Volume FOUND: ", name)
				return volume, true, nil
			}
		}
	}
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Volume NOT Found: ", name)
	return nil, false, nil
}

// List returns array of types.Containers and error
func List(docker *client.Docker, all bool) (types.VolumesListResponse, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Docker.Timeout*time.Second)
	// defer cancel()
	filter := filters.Args{}
	volumes, err := docker.Client.VolumeList(context.Background(), filter)
	if err != nil {
		return types.VolumesListResponse{}, err
	}
	return volumes, nil
}
