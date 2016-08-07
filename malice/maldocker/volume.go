package maldocker

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/strslice"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/maliceio/malice/malice/persist"
	"golang.org/x/net/context"
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

func resolveLocalPath(localPath string) (absPath string, err error) {
	if absPath, err = filepath.Abs(localPath); err != nil {
		return
	}

	return archive.PreserveTrailingDotOrSeparator(absPath, localPath), nil
}

func (client *Docker) statContainerPath(containerName, path string) (types.ContainerPathStat, error) {
	return client.Client.ContainerStatPath(context.Background(), containerName, path)
}

// CopyToVolume copies samples into Malice volume
func (client *Docker) CopyToVolume(file persist.File) {

	name := "copy2volume"
	image := "busybox"
	cmd := strslice.StrSlice{"sh", "-c", "while true; do echo 'Waiting...'; sleep 1; done"}
	binds := []string{"malice:/malice:rw"}
	volSavePath := filepath.Join("/malice", file.SHA256)

	if client.Ping() {
		container, err := client.StartContainer(cmd, name, image, false, binds, nil, nil, nil)
		er.CheckError(err)

		// Prepare destination copy info by stat-ing the container path.
		dstInfo := archive.CopyInfo{Path: volSavePath}
		dstStat, err := client.statContainerPath(container.Name, volSavePath)
		log.WithFields(log.Fields{
			"dstInfo":        dstInfo,
			"dstStat":        dstStat,
			"container.Name": container.Name,
			"file.Path":      file.Path,
			"volSavePath":    volSavePath,
			"SampledsDir":    maldirs.GetSampledsDir(),
		}).Debug("First statContainerPath call.")
		// er.CheckError(err)

		// Check if file already exists in volume
		if dstStat.Size > 0 {
			// Remove copy2volume container
			client.RemoveContainer(container, true, true, true)
			log.Debug("Sample ", file.Name, " already in malice volume.")
			return
		}

		// If the destination is a symbolic link, we should evaluate it.
		if err == nil && dstStat.Mode&os.ModeSymlink != 0 {
			linkTarget := dstStat.LinkTarget
			if !system.IsAbs(linkTarget) {
				// Join with the parent directory.
				dstParent, _ := archive.SplitPathDirEntry(volSavePath)
				linkTarget = filepath.Join(dstParent, linkTarget)
			}

			dstInfo.Path = linkTarget
			dstStat, err = client.statContainerPath(container.Name, linkTarget)
			log.WithFields(log.Fields{
				"dstInfo":        dstInfo,
				"dstStat":        dstStat,
				"container.Name": container.Name,
				"file.Path":      file.Path,
				"linkTarget":     linkTarget,
				"SampledsDir":    maldirs.GetSampledsDir(),
			}).Debug("Second statContainerPath call.")
			er.CheckError(err)
		}

		if err == nil {
			dstInfo.Exists, dstInfo.IsDir = true, dstStat.Mode.IsDir()
		}

		var (
			content         io.Reader
			resolvedDstPath string
		)

		// Prepare source copy info.
		srcInfo, err := archive.CopyInfoSourcePath(file.Path, false)
		er.CheckError(err)

		srcArchive, err := archive.TarResource(srcInfo)
		er.CheckError(err)
		defer srcArchive.Close()

		dstDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
		er.CheckError(err)

		defer preparedArchive.Close()

		resolvedDstPath = dstDir
		content = preparedArchive

		copyOptions := types.CopyToContainerOptions{
			AllowOverwriteDirWithFile: true,
		}

		// Copy sample to malice volume
		er.CheckError(client.Client.CopyToContainer(
			context.Background(),
			container.ID,
			resolvedDstPath,
			content,
			copyOptions,
		))

		// Remove copy2volume container
		client.RemoveContainer(container, true, true, true)
	}
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
