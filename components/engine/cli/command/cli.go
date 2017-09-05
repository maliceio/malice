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
	cliconfig "github.com/maliceio/engine/cli/config"
	"github.com/maliceio/engine/cli/config/configfile"
	cliflags "github.com/maliceio/engine/cli/flags"
	"github.com/maliceio/engine/client"
	"github.com/maliceio/engine/malice/version"
	mopts "github.com/maliceio/engine/opts"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Cli represents the malice command line client.
type Cli interface {
	Client() client.APIClient
	Err() io.Writer
	ConfigFile() *configfile.ConfigFile
	ServerInfo() ServerInfo
}

// MaliceCli is an instance the malice command line client.
// Instances of the client can be returned from NewMaliceCli.
type MaliceCli struct {
	configFile     *configfile.ConfigFile
	err            io.Writer
	client         client.APIClient
	defaultVersion string
	server         ServerInfo
}

// DefaultVersion returns api.defaultVersion or DOCKER_API_VERSION if specified.
func (cli *MaliceCli) DefaultVersion() string {
	return cli.defaultVersion
}

// Client returns the APIClient
func (cli *MaliceCli) Client() client.APIClient {
	return cli.client
}

// Err returns the writer used for stderr
func (cli *MaliceCli) Err() io.Writer {
	return cli.err
}

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOutput(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}

// ConfigFile returns the ConfigFile
func (cli *MaliceCli) ConfigFile() *configfile.ConfigFile {
	return cli.configFile
}

// ServerInfo returns the server version details for the host this client is
// connected to
func (cli *MaliceCli) ServerInfo() ServerInfo {
	return cli.server
}

// Initialize the maliceCli runs initialization that must happen after command
// line flags are parsed.
func (cli *MaliceCli) Initialize(opts *cliflags.ClientOptions) error {
	cli.configFile = LoadDefaultConfigFile(cli.err)

	var err error
	cli.client, err = NewAPIClientFromFlags(opts.Common, cli.configFile)
	// if tlsconfig.IsErrEncryptedKey(err) {
	// 	var (
	// 		passwd string
	// 		giveup bool
	// 	)
	// 	passRetriever := passphrase.PromptRetrieverWithInOut(cli.In(), cli.Out(), nil)

	// 	for attempts := 0; tlsconfig.IsErrEncryptedKey(err); attempts++ {
	// 		// some code and comments borrowed from notary/trustmanager/keystore.go
	// 		passwd, giveup, err = passRetriever("private", "encrypted TLS private", false, attempts)
	// 		// Check if the passphrase retriever got an error or if it is telling us to give up
	// 		if giveup || err != nil {
	// 			return errors.Wrap(err, "private key is encrypted, but could not get passphrase")
	// 		}

	// 		opts.Common.TLSOptions.Passphrase = passwd
	// 		cli.client, err = NewAPIClientFromFlags(opts.Common, cli.configFile)
	// 	}
	// }

	if err != nil {
		return err
	}

	cli.defaultVersion = cli.client.ClientVersion()

	if ping, err := cli.client.Ping(context.Background()); err == nil {
		cli.server = ServerInfo{
			OSType: ping.OSType,
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

// NewMaliceCli returns a MaliceCli instance with IO output and error streams set by in, out and err.
func NewMaliceCli(err io.Writer) *MaliceCli {
	return &MaliceCli{err: err}
}

// LoadDefaultConfigFile attempts to load the default config file and returns
// an initialized ConfigFile struct if none is found.
func LoadDefaultConfigFile(err io.Writer) *configfile.ConfigFile {
	configFile, e := cliconfig.Load(cliconfig.Dir())
	if e != nil {
		fmt.Fprintf(err, "WARNING: Error loading config file:%v\n", e)
	}
	// if !configFile.ContainsAuth() {
	// 	credentials.DetectDefaultStore(configFile)
	// }
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
	if tmpStr := os.Getenv("MALICE_API_VERSION"); tmpStr != "" {
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
		host = os.Getenv("MALICE_HOST")
	case 1:
		host = hosts[0]
	default:
		return "", errors.New("Please specify only one -H")
	}

	host, err = mopts.ParseHost(tlsOptions != nil, host)
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
	return "Malice-Client/" + version.Version + " (" + runtime.GOOS + ")"
}
