package container

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/maliceio/malice/malice/docker/client"
	"golang.org/x/net/context"
)

// Remove removes the `cont` container unforcedly.
// If volumes is true, the associated volumes are removed with container.
// If links is true, the associated links are removed with container.
// If force is true, the container will be destroyed with extreme prejudice.
func Remove(docker *client.Docker, contID string, volumes bool, links bool, force bool) error {
	log.Debug("Removing container: ", contID)
	return removeContainer(docker, context.Background(), contID, volumes, links, force)
	// // check if container exists
	// if plugin, exists, _ := Exists(docker, cont.Name); exists {
	// 	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Removing Plugin container: ", cont.Name)
	// 	er.CheckError(docker.Client.ContainerRemove(context.Background(), plugin.ID, types.ContainerRemoveOptions{
	// 		RemoveVolumes: true,
	// 		// RemoveLinks:   links,
	// 		Force: true,
	// 	}))
	// } else {
	// 	// container not found
	// 	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Error("Plugin container does not exist. Cannot remove.")
	// }
	// return nil
}

// removeContainer
func removeContainer(docker *client.Docker, ctx context.Context, container string, removeVolumes, removeLinks, force bool) error {
	// name = strings.Trim(name, "/")
	options := types.ContainerRemoveOptions{
		RemoveVolumes: removeVolumes,
		RemoveLinks:   removeLinks,
		Force:         force,
	}
	if err := docker.Client.ContainerRemove(ctx, container, options); err != nil {
		return err
	}
	return nil
}
