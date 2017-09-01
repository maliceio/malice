package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/maliceio/engine/api"
	"github.com/maliceio/engine/api/types/versions"
	"github.com/maliceio/engine/cli/config/configfile"
	"github.com/maliceio/engine/client"
	"github.com/pkg/errors"
)

// Cli represents the docker command line client.
type Cli interface {
	Client() client.APIClient
	Err() io.Writer
	ConfigFile() *configfile.ConfigFile
}

// DockerCli is an instance the docker command line client.
// Instances of the client can be returned from NewDockerCli.
type DockerCli struct {
	configFile     *configfile.ConfigFile
	err            io.Writer
	client         client.APIClient
	defaultVersion string
	server         ServerInfo
}

// Initialize the dockerCli runs initialization that must happen after command
// line flags are parsed.
func (cli *DockerCli) Initialize(opts *cliflags.ClientOptions) error {
	cli.configFile = LoadDefaultConfigFile(cli.err)

	var err error
	cli.client, err = NewAPIClientFromFlags(opts.Common, cli.configFile)
	if tlsconfig.IsErrEncryptedKey(err) {
		var (
			passwd string
			giveup bool
		)
		passRetriever := passphrase.PromptRetrieverWithInOut(cli.In(), cli.Out(), nil)

		for attempts := 0; tlsconfig.IsErrEncryptedKey(err); attempts++ {
			// some code and comments borrowed from notary/trustmanager/keystore.go
			passwd, giveup, err = passRetriever("private", "encrypted TLS private", false, attempts)
			// Check if the passphrase retriever got an error or if it is telling us to give up
			if giveup || err != nil {
				return errors.Wrap(err, "private key is encrypted, but could not get passphrase")
			}

			opts.Common.TLSOptions.Passphrase = passwd
			cli.client, err = NewAPIClientFromFlags(opts.Common, cli.configFile)
		}
	}

	if err != nil {
		return err
	}

	cli.defaultVersion = cli.client.ClientVersion()

	if ping, err := cli.client.Ping(context.Background()); err == nil {
		cli.server = ServerInfo{
			HasExperimental: ping.Experimental,
			OSType:          ping.OSType,
		}

		// since the new header was added in 1.25, assume server is 1.24 if header is not present.
		if ping.APIVersion == "" {
			ping.APIVersion = "1.24"
		}

		// if server version is lower than the current cli, downgrade
		if versions.LessThan(ping.APIVersion, cli.client.ClientVersion()) {
			cli.client.UpdateClientVersion(ping.APIVersion)
		}
	}

	return nil
}

// ServerInfo stores details about the supported features and platform of the
// server
type ServerInfo struct {
	HasExperimental bool
	OSType          string
}

// NewDockerCli returns a DockerCli instance with IO output and error streams set by in, out and err.
func NewDockerCli(in io.ReadCloser, out, err io.Writer) *DockerCli {
	return &DockerCli{in: NewInStream(in), out: NewOutStream(out), err: err}
}

// LoadDefaultConfigFile attempts to load the default config file and returns
// an initialized ConfigFile struct if none is found.
func LoadDefaultConfigFile(err io.Writer) *configfile.ConfigFile {
	configFile, e := cliconfig.Load(cliconfig.Dir())
	if e != nil {
		fmt.Fprintf(err, "WARNING: Error loading config file:%v\n", e)
	}
	if !configFile.ContainsAuth() {
		credentials.DetectDefaultStore(configFile)
	}
	return configFile
}

// NewAPIClientFromFlags creates a new APIClient from command line flags
func NewAPIClientFromFlags(opts *cliflags.CommonOptions, configFile *configfile.ConfigFile) (client.APIClient, error) {
	host, err := getServerHost(opts.Hosts, opts.TLSOptions)
	if err != nil {
		return &client.Client{}, err
	}

	customHeaders := configFile.HTTPHeaders
	if customHeaders == nil {
		customHeaders = map[string]string{}
	}
	customHeaders["User-Agent"] = UserAgent()

	verStr := api.DefaultVersion
	if tmpStr := os.Getenv("DOCKER_API_VERSION"); tmpStr != "" {
		verStr = tmpStr
	}

	httpClient, err := newHTTPClient(host, opts.TLSOptions)
	if err != nil {
		return &client.Client{}, err
	}

	return client.NewClient(host, verStr, httpClient, customHeaders)
}

func getServerHost(hosts []string, tlsOptions *tlsconfig.Options) (host string, err error) {
	switch len(hosts) {
	case 0:
		host = os.Getenv("DOCKER_HOST")
	case 1:
		host = hosts[0]
	default:
		return "", errors.New("Please specify only one -H")
	}

	host, err = dopts.ParseHost(tlsOptions != nil, host)
	return
}

func newHTTPClient(host string, tlsOptions *tlsconfig.Options) (*http.Client, error) {
	if tlsOptions == nil {
		// let the api client configure the default transport.
		return nil, nil
	}
	opts := *tlsOptions
	opts.ExclusiveRootPools = true
	config, err := tlsconfig.Client(opts)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: config,
	}
	proto, addr, _, err := client.ParseHost(host)
	if err != nil {
		return nil, err
	}

	sockets.ConfigureTransport(tr, proto, addr)

	return &http.Client{
		Transport: tr,
	}, nil
}

// UserAgent returns the user agent string used for making API requests
func UserAgent() string {
	return "Docker-Client/" + cli.Version + " (" + runtime.GOOS + ")"
}
