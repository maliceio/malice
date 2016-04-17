package maldocker

import (
	"time"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/maliceio/malice/config"
	"golang.org/x/net/context"

	"regexp"

	log "github.com/Sirupsen/logrus"
)

// VolumeExists returns type.Volume and true
// if the volume name exists, otherwise false.
func (client *Docker) VolumeExists(name string) (*types.Volume, bool, error) {
	return client.ParseVolumes(name, true)
}

// CreateVolume creates a docker volume with the given name
// returns: Volume, error
func (client *Docker) CreateVolume(name string) (types.Volume, error) {
	options := types.VolumeCreateRequest{
		Name: name, // Name is the requested name of the volume
		// Driver     string            // Driver is the name of the driver that should be used to create the volume
		// DriverOpts map[string]string // DriverOpts holds the driver specific options to use for when creating the volume.
		// Labels     map[string]string // Labels holds metadata specific to the volume being created.
	}
	vol, err := client.Client.VolumeCreate(context.Background(), options)
	log.WithFields(log.Fields{
		"name": name,
		"env":  config.Conf.Environment.Run,
	}).Info("Created Volume: ", name)
	return vol, err
}

// ParseVolumes parses the volumes
func (client *Docker) ParseVolumes(name string, all bool) (*types.Volume, bool, error) {
	// list volumes
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Searching for volume: ", name)
	volumes, err := client.listVolumes(all)
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

// listVolumes returns array of types.Containers and error
func (client *Docker) listVolumes(all bool) (types.VolumesListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Docker.Timeout*time.Second)
	defer cancel()
	filter := filters.Args{}
	volumes, err := client.Client.VolumeList(ctx, filter)
	if err != nil {
		return types.VolumesListResponse{}, err
	}
	return volumes, nil
}
