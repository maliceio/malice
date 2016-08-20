package image

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"golang.org/x/net/context"

	"regexp"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/dockerignore"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/stringutils"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/docker/docker/registry"
	"github.com/docker/engine-api/types"
	registrytypes "github.com/docker/engine-api/types/registry"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/docker/client"
	er "github.com/maliceio/malice/malice/errors"

	log "github.com/Sirupsen/logrus"
)

// Pull pulls docker image:tag
// TODO: add trusted pull for offcial malice plugins
func Pull(docker *client.Docker, id string, tag string) {

	responseBody, err := docker.Client.ImagePull(context.Background(), id, types.ImagePullOptions{})
	defer responseBody.Close()
	er.CheckError(err)

	jsonmessage.DisplayJSONMessagesStream(responseBody, os.Stdout, os.Stdout.Fd(), true, nil)
}

// Build builds docker image from git repository
func Build(docker *client.Docker, repository string, tags []string, buildArgs map[string]string, labels map[string]string, quiet bool) {

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
	response, err := docker.Client.ImageBuild(context.Background(), body, buildOptions)
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

// Exists returns APIImages images list and true
// if the image name exists, otherwise false.
func Exists(docker *client.Docker, name string) (types.Image, bool, error) {
	log.WithFields(log.Fields{"env": config.Conf.Environment.Run}).Debug("Searching for image: ", name)
	images, err := List(docker, name, false)
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

// List lists all images
func List(docker *client.Docker, name string, all bool) ([]types.Image, error) {

	options := types.ImageListOptions{
		All:       all,
		MatchName: name,
		// Filters   filters.Args
	}
	imageList, err := docker.Client.ImageList(context.Background(), options)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return imageList, nil
}

type searchOptions struct {
	term    string
	noTrunc bool
	limit   int
	filter  []string

	// Deprecated
	stars     uint
	automated bool
}

// Search searches for malice images
func Search(docker *client.Docker, term string) error {

	opts := searchOptions{
		term:  term,
		limit: registry.DefaultSearchLimit,
	}
	options := types.ImageSearchOptions{
		// RegistryAuth:  encodedAuth,
		// PrivilegeFunc: requestPrivilege,
		// Filters: searchFilters,
		Limit: opts.limit,
	}

	unorderedResults, err := docker.Client.ImageSearch(context.Background(), opts.term, options)
	if err != nil {
		return err
	}

	results := searchResultsByStars(unorderedResults)
	sort.Sort(results)

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
	fmt.Fprintf(w, "NAME\tDESCRIPTION\tSTARS\tOFFICIAL\tAUTOMATED\n")
	for _, res := range results {
		// --automated and -s, --stars are deprecated since Docker 1.12
		if (opts.automated && !res.IsAutomated) || (int(opts.stars) > res.StarCount) {
			continue
		}
		desc := strings.Replace(res.Description, "\n", " ", -1)
		desc = strings.Replace(desc, "\r", " ", -1)
		if !opts.noTrunc && len(desc) > 45 {
			desc = stringutils.Truncate(desc, 42) + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t", res.Name, desc, res.StarCount)
		if res.IsOfficial {
			fmt.Fprint(w, "[OK]")

		}
		fmt.Fprint(w, "\t")
		if res.IsAutomated {
			fmt.Fprint(w, "[OK]")
		}
		fmt.Fprint(w, "\n")
	}
	w.Flush()
	return nil
}

// SearchResultsByStars sorts search results in descending order by number of stars.
type searchResultsByStars []registrytypes.SearchResult

func (r searchResultsByStars) Len() int           { return len(r) }
func (r searchResultsByStars) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r searchResultsByStars) Less(i, j int) bool { return r[j].StarCount < r[i].StarCount }
