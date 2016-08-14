package maldocker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/dockerignore"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/docker/engine-api/types"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"

	"regexp"

	log "github.com/Sirupsen/logrus"
)

// PullImage pulls docker image:tag
func (client *Docker) PullImage(id string, tag string) {

	responseBody, err := client.Client.ImagePull(context.Background(), id, types.ImagePullOptions{})
	defer responseBody.Close()
	er.CheckError(err)

	jsonmessage.DisplayJSONMessagesStream(responseBody, os.Stdout, os.Stdout.Fd(), true, nil)
}

// BuildImage builds docker image from git repository
func (client *Docker) BuildImage(repository string, tags []string, buildArgs map[string]string, labels map[string]string, quiet bool) {

	var (
		buildCtx io.ReadCloser
		err      error
	)

	var (
		contextDir    string
		tempDir       string
		relDockerfile string
		progBuff      io.Writer
		buildBuff     io.Writer
	)

	progBuff = os.Stdout
	buildBuff = os.Stdout

	switch {
	case repository == "-":
		buildCtx, relDockerfile, err = builder.GetContextFromReader(os.Stdin, "")
	case urlutil.IsGitURL(repository):
		tempDir, relDockerfile, err = builder.GetContextFromGitURL(repository, "")
	case urlutil.IsURL(repository):
		buildCtx, relDockerfile, err = builder.GetContextFromURL(progBuff, repository, "")
	default:
		_, relDockerfile, err = builder.GetContextFromLocalDir(repository, "")
	}

	if tempDir != "" {
		defer os.RemoveAll(tempDir)
		contextDir = tempDir
	}

	if buildCtx == nil {
		// And canonicalize dockerfile name to a platform-independent one
		relDockerfile, err = archive.CanonicalTarNameForPath(relDockerfile)
		if err != nil {
			log.Fatalf("cannot canonicalize dockerfile path %s: %v", relDockerfile, err)
		}

		f, err := os.Open(filepath.Join(contextDir, ".dockerignore"))
		if err != nil && !os.IsNotExist(err) {
			er.CheckError(err)
		}
		defer f.Close()

		var excludes []string
		if err == nil {
			excludes, err = dockerignore.ReadAll(f)
			if err != nil {
				er.CheckError(err)
			}
		}

		if err := builder.ValidateContextDirectory(contextDir, excludes); err != nil {
			log.Fatalf("Error checking context: '%s'.", err)
		}

		// If .dockerignore mentions .dockerignore or the Dockerfile
		// then make sure we send both files over to the daemon
		// because Dockerfile is, obviously, needed no matter what, and
		// .dockerignore is needed to know if either one needs to be
		// removed. The daemon will remove them for us, if needed, after it
		// parses the Dockerfile. Ignore errors here, as they will have been
		// caught by validateContextDirectory above.
		var includes = []string{"."}
		keepThem1, _ := fileutils.Matches(".dockerignore", excludes)
		keepThem2, _ := fileutils.Matches(relDockerfile, excludes)
		if keepThem1 || keepThem2 {
			includes = append(includes, ".dockerignore", relDockerfile)
		}

		buildCtx, err = archive.TarWithOptions(contextDir, &archive.TarOptions{
			Compression:     archive.Uncompressed,
			ExcludePatterns: excludes,
			IncludeFiles:    includes,
		})
		er.CheckError(err)
	}

	// Setup an upload progress bar
	progressOutput := streamformatter.NewStreamFormatter().NewProgressOutput(progBuff, true)

	var body io.Reader = progress.NewProgressReader(buildCtx, progressOutput, 0, "", "Sending build context to Docker daemon")

	buildOptions := types.ImageBuildOptions{
		Tags:           tags,
		SuppressOutput: quiet,
		// RemoteContext  string
		NoCache: true,
		// Remove         bool
		// ForceRemove    bool
		// PullParent     bool
		// Isolation      container.Isolation
		// CPUSetCPUs     string
		// CPUSetMems     string
		// CPUShares      int64
		// CPUQuota       int64
		// CPUPeriod      int64
		// Memory         int64
		// MemorySwap     int64
		// CgroupParent   string
		// ShmSize        int64
		Dockerfile: relDockerfile,
		// Ulimits        []*units.Ulimit
		BuildArgs: buildArgs,
		// AuthConfigs    map[string]AuthConfig
		// Context        io.Reader
		Labels: labels,
	}
	response, err := client.Client.ImageBuild(context.Background(), body, buildOptions)
	defer response.Body.Close()
	er.CheckError(err)

	err = jsonmessage.DisplayJSONMessagesStream(response.Body, buildBuff, os.Stdout.Fd(), true, nil)
	if err != nil {
		if jerr, ok := err.(*jsonmessage.JSONError); ok {
			// If no error code is set, default to 1
			if jerr.Code == 0 {
				jerr.Code = 1
			}
			if quiet {
				fmt.Fprintf(os.Stderr, "%s%s", progBuff, buildBuff)
			}
		}
	}
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
	// ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Docker.Timeout*time.Second)
	// defer cancel()

	options := types.ImageListOptions{
		All: all,
		// MatchName string
		// Filters   filters.Args
	}
	imageList, err := client.Client.ImageList(context.Background(), options)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return imageList, nil
}
