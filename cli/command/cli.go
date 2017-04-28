package command

import (
	"fmt"
	"io"

	"github.com/moby/moby/api/types"
	cliconfig "github.com/moby/moby/cli/config"
	"github.com/moby/moby/cli/config/configfile"
	"github.com/moby/moby/cli/config/credentials"
	cliflags "github.com/moby/moby/cli/flags"
	"github.com/moby/moby/client"
	"github.com/spf13/cobra"
)

// Streams is an interface which exposes the standard input and output streams
type Streams interface {
	// In() *InStream
	// Out() *OutStream
	// Err() io.Writer
}

// Cli represents the docker command line client.
type Cli interface {
	Client() client.APIClient
	// Out() *OutStream
	// Err() io.Writer
	// In() *InStream
	ConfigFile() *configfile.ConfigFile
}

// MaliceCli is an instance the docker command line client.
// Instances of the client can be returned from NewMaliceCli.
type MaliceCli struct {
	configFile     *configfile.ConfigFile
	in             io.ReadCloser
	out            io.Writer
	err            io.Writer
	keyFile        string
	client         client.APIClient
	defaultVersion string
	// server         ServerInfo
}

// DefaultVersion returns api.defaultVersion or DOCKER_API_VERSION if specified.
func (cli *MaliceCli) DefaultVersion() string {
	return cli.defaultVersion
}

// Client returns the APIClient
func (cli *MaliceCli) Client() client.APIClient {
	return cli.client
}

// Out returns the writer used for stdout
func (cli *MaliceCli) Out() io.Writer {
	return cli.out
}

// Err returns the writer used for stderr
func (cli *MaliceCli) Err() io.Writer {
	return cli.err
}

// In returns the reader used for stdin
func (cli *MaliceCli) In() io.ReadCloser {
	return cli.in
}

// ShowHelp shows the command help.
func (cli *MaliceCli) ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.SetOutput(cli.err)
	cmd.HelpFunc()(cmd, args)
	return nil
}

// ConfigFile returns the ConfigFile
func (cli *MaliceCli) ConfigFile() *configfile.ConfigFile {
	return cli.configFile
}

func addAll(to, from map[string]types.AuthConfig) {
	for reg, ac := range from {
		to[reg] = ac
	}
}

// Initialize the MaliceCli runs initialization that must happen after command
// line flags are parsed.
func (cli *MaliceCli) Initialize(opts *cliflags.ClientOptions) error {
	cli.configFile = LoadDefaultConfigFile(cli.err)

	var err error

	//===================================================================================================
	// TODO: I should create my own API and have my cli comsume it as well, damn those docker peeps are legends! //////////////////////////////////////////////////
	//===================================================================================================

	// cli.client, err = NewAPIClientFromFlags(opts.Common, cli.configFile)
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

	// cli.defaultVersion = cli.client.ClientVersion()

	//===================================================================================================
	// TODO: I should try to ping the docker cli here and/or error out. //////////////////////////////////////////////////////////////////////////////////
	//===================================================================================================

	// if ping, err := cli.client.Ping(context.Background()); err == nil {
	// 	cli.server = ServerInfo{
	// 		HasExperimental: ping.Experimental,
	// 		OSType:          ping.OSType,
	// 	}

	// 	// since the new header was added in 1.25, assume server is 1.24 if header is not present.
	// 	if ping.APIVersion == "" {
	// 		ping.APIVersion = "1.24"
	// 	}

	// 	// if server version is lower than the current cli, downgrade
	// 	if versions.LessThan(ping.APIVersion, cli.client.ClientVersion()) {
	// 		cli.client.UpdateClientVersion(ping.APIVersion)
	// 	}
	// }

	return nil
}

// NewMaliceCli returns a MaliceCli instance with IO output and error streams set by in, out and err.
func NewMaliceCli(in io.ReadCloser, out, err io.Writer) *MaliceCli {
	return &MaliceCli{in: in, out: out, err: err}
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

// // NewAPIClientFromFlags creates a new APIClient from command line flags
// func NewAPIClientFromFlags(opts *cliflags.CommonOptions, configFile *configfile.ConfigFile) (client.APIClient, error) {
// 	host, err := getServerHost(opts.Hosts, opts.TLSOptions)
// 	if err != nil {
// 		return &client.Client{}, err
// 	}

// 	customHeaders := configFile.HTTPHeaders
// 	if customHeaders == nil {
// 		customHeaders = map[string]string{}
// 	}
// 	customHeaders["User-Agent"] = UserAgent()

// 	verStr := api.DefaultVersion
// 	if tmpStr := os.Getenv("DOCKER_API_VERSION"); tmpStr != "" {
// 		verStr = tmpStr
// 	}

// 	httpClient, err := newHTTPClient(host, opts.TLSOptions)
// 	if err != nil {
// 		return &client.Client{}, err
// 	}

// 	return client.NewClient(host, verStr, httpClient, customHeaders)
// }

// func getServerHost(hosts []string, tlsOptions *tlsconfig.Options) (host string, err error) {
// 	switch len(hosts) {
// 	case 0:
// 		host = os.Getenv("DOCKER_HOST")
// 	case 1:
// 		host = hosts[0]
// 	default:
// 		return "", errors.New("Please specify only one -H")
// 	}

// 	host, err = dopts.ParseHost(tlsOptions != nil, host)
// 	return
// }

// func newHTTPClient(host string, tlsOptions *tlsconfig.Options) (*http.Client, error) {
// 	if tlsOptions == nil {
// 		// let the api client configure the default transport.
// 		return nil, nil
// 	}
// 	opts := *tlsOptions
// 	opts.ExclusiveRootPools = true
// 	config, err := tlsconfig.Client(opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tr := &http.Transport{
// 		TLSClientConfig: config,
// 	}
// 	proto, addr, _, err := client.ParseHost(host)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sockets.ConfigureTransport(tr, proto, addr)

// 	return &http.Client{
// 		Transport: tr,
// 	}, nil
// }

// // UserAgent returns the user agent string used for making API requests
// func UserAgent() string {
// 	return "Docker-Client/" + dockerversion.Version + " (" + runtime.GOOS + ")"
// }
