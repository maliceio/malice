package container

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/system"
	"github.com/maliceio/malice/malice/docker/client"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
	"github.com/maliceio/malice/malice/persist"
)

// CopyToVolume copies samples into Malice volume
func CopyToVolume(docker *client.Docker, file persist.File) {

	name := "copy2volume"
	image := "busybox"
	cmd := strslice.StrSlice{"sh", "-c", "while true; do echo 'Waiting...'; sleep 1; done"}
	binds := []string{"malice:/malice:rw"}
	volSavePath := filepath.Join("/malice", file.SHA256)

	if docker.Ping() {
		cont, err := Start(docker, cmd, name, image, false, binds, nil, nil, nil)
		er.CheckError(err)

		defer func() {
			er.CheckError(Remove(docker, cont.ID, true, false, true))
		}()

		// Get an absolute source path.
		srcPath, err := resolveLocalPath(file.Path)
		er.CheckError(err)

		// Prepare destination copy info by stat-ing the container path.
		dstInfo := archive.CopyInfo{Path: volSavePath}
		dstStat, err := statContainerPath(docker, cont.Name, volSavePath)
		log.WithFields(log.Fields{
			"dstInfo":        dstInfo,
			"dstStat":        dstStat,
			"container.Name": cont.Name,
			"file.Path":      file.Path,
			"volSavePath":    volSavePath,
			"SampledsDir":    maldirs.GetSampledsDir(),
		}).Debug("First statContainerPath call.")
		// er.CheckError(err)

		// Check if file already exists in volume
		if dstStat.Size > 0 {
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
			dstStat, err = statContainerPath(docker, cont.Name, linkTarget)
			log.WithFields(log.Fields{
				"dstInfo":        dstInfo,
				"dstStat":        dstStat,
				"container.Name": cont.Name,
				"file.Path":      file.Path,
				"linkTarget":     linkTarget,
				"SampledsDir":    maldirs.GetSampledsDir(),
			}).Debug("Second statContainerPath call.")
			er.CheckError(err)
		}

		if err == nil {
			dstInfo.Exists = true
			dstInfo.IsDir = dstStat.Mode.IsDir()
		}

		var (
			content         io.Reader
			resolvedDstPath string
		)

		// Prepare source copy info.
		srcInfo, err := archive.CopyInfoSourcePath(srcPath, false)
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
		er.CheckError(docker.Client.CopyToContainer(
			context.Background(),
			cont.ID,
			resolvedDstPath,
			content,
			copyOptions,
		))
	}
}

func resolveLocalPath(localPath string) (absPath string, err error) {
	if absPath, err = filepath.Abs(localPath); err != nil {
		return
	}

	return archive.PreserveTrailingDotOrSeparator(absPath, localPath), nil
}

func statContainerPath(docker *client.Docker, containerName, path string) (types.ContainerPathStat, error) {
	return docker.Client.ContainerStatPath(context.Background(), containerName, path)
}
