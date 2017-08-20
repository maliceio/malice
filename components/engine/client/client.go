/*
Package client is a Go client for the Malice Engine API.

The "malice" command uses this package to communicate with the daemon. It can also
be used by your own Go applications to do anything the command-line interface does
- running containers, pulling images, managing swarms, etc.

For more information about the Engine API, see the documentation:
https://docs.malice.io/engine/reference/api/

Usage

You use the library by creating a client object and calling methods on it. The
client can be created either from environment variables with NewEnvClient, or
configured manually with NewClient.

For example, to list running containers (the equivalent of "malice plugin ls"):

	package main

	import (
		"context"
		"fmt"

		"github.com/maliceio/engine/api/types"
		"github.com/maliceio/engine/client"
	)

	func main() {
		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		plugins, err := cli.PluginList(context.Background(), types.PluginListListOptions{})
		if err != nil {
			panic(err)
		}

		for _, plugin := range plugins {
			fmt.Printf("%s %s\n", plugin.Name, plugin.Image)
		}
	}

*/
package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/maliceio/engine/api"
	"github.com/maliceio/engine/api/types"
	"github.com/maliceio/engine/api/types/versions"
	"golang.org/x/net/context"
)

// DefaultMaliceHost defines os specific default if MALICE_HOST is unset
const DefaultMaliceHost = "unix:///var/run/malice.sock"

// ErrRedirect is the error returned by checkRedirect when the request is non-GET.
var ErrRedirect = errors.New("unexpected redirect in response")

// Client is the API client that performs all operations
// against a docker server.
type Client struct {
	// scheme sets the scheme for the client
	scheme string
	// host holds the server address to connect to
	host string
	// proto holds the client protocol i.e. unix.
	proto string
	// addr holds the client address.
	addr string
	// basePath holds the path to prepend to the requests.
	basePath string
	// client used to send and receive http requests.
	client *http.Client
	// version of the server to talk to.
	version string
	// custom http headers configured by users.
	customHTTPHeaders map[string]string
	// manualOverride is set to true when the version was set by users.
	manualOverride bool
}

// CheckRedirect specifies the policy for dealing with redirect responses:
// If the request is non-GET return `ErrRedirect`. Otherwise use the last response.
//
// Go 1.8 changes behavior for HTTP redirects (specifically 301, 307, and 308) in the client .
// The Malice client (and by extension docker API client) can be made to to send a request
// like POST /containers//start where what would normally be in the name section of the URL is empty.
// This triggers an HTTP 301 from the daemon.
// In go 1.8 this 301 will be converted to a GET request, and ends up getting a 404 from the daemon.
// This behavior change manifests in the client in that before the 301 was not followed and
// the client did not generate an error, but now results in a message like Error response from daemon: page not found.
func CheckRedirect(req *http.Request, via []*http.Request) error {
	if via[0].Method == http.MethodGet {
		return http.ErrUseLastResponse
	}
	return ErrRedirect
}

// NewEnvClient initializes a new API client based on environment variables.
// Use MALICE_HOST to set the url to the docker server.
// Use MALICE_API_VERSION to set the version of the API to reach, leave empty for latest.
// Use MALICE_CERT_PATH to load the TLS certificates from.
// Use MALICE_TLS_VERIFY to enable or disable TLS verification, off by default.
func NewEnvClient() (*Client, error) {
	var client *http.Client
	if dockerCertPath := os.Getenv("MALICE_CERT_PATH"); dockerCertPath != "" {
		options := tlsconfig.Options{
			CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
			CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
			KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
			InsecureSkipVerify: os.Getenv("MALICE_TLS_VERIFY") == "",
		}
		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			return nil, err
		}

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
			CheckRedirect: CheckRedirect,
		}
	}

	host := os.Getenv("MALICE_HOST")
	if host == "" {
		host = DefaultMaliceHost
	}
	version := os.Getenv("MALICE_API_VERSION")
	if version == "" {
		version = api.DefaultVersion
	}

	cli, err := NewClient(host, version, client, nil)
	if err != nil {
		return cli, err
	}
	if os.Getenv("MALICE_API_VERSION") != "" {
		cli.manualOverride = true
	}
	return cli, nil
}

// NewClient initializes a new API client for the given host and API version.
// It uses the given http client as transport.
// It also initializes the custom http headers to add to each request.
//
// It won't send any version information if the version number is empty. It is
// highly recommended that you set a version or your client may break if the
// server is upgraded.
func NewClient(host string, version string, client *http.Client, httpHeaders map[string]string) (*Client, error) {
	proto, addr, basePath, err := ParseHost(host)
	if err != nil {
		return nil, err
	}

	if client != nil {
		if _, ok := client.Transport.(http.RoundTripper); !ok {
			return nil, fmt.Errorf("unable to verify TLS configuration, invalid transport %v", client.Transport)
		}
	} else {
		transport := new(http.Transport)
		sockets.ConfigureTransport(transport, proto, addr)
		client = &http.Client{
			Transport:     transport,
			CheckRedirect: CheckRedirect,
		}
	}

	scheme := "http"
	tlsConfig := resolveTLSConfig(client.Transport)
	if tlsConfig != nil {
		scheme = "https"
	}

	return &Client{
		scheme:            scheme,
		host:              host,
		proto:             proto,
		addr:              addr,
		basePath:          basePath,
		client:            client,
		version:           version,
		customHTTPHeaders: httpHeaders,
	}, nil
}

// Close ensures that transport.Client is closed
// especially needed while using NewClient with *http.Client = nil
// for example
// client.NewClient("unix:///var/run/docker.sock", nil, "v1.18", map[string]string{"User-Agent": "engine-api-cli-1.0"})
func (cli *Client) Close() error {

	if t, ok := cli.client.Transport.(*http.Transport); ok {
		t.CloseIdleConnections()
	}

	return nil
}

// getAPIPath returns the versioned request path to call the api.
// It appends the query parameters to the path if they are not empty.
func (cli *Client) getAPIPath(p string, query url.Values) string {
	var apiPath string
	if cli.version != "" {
		v := strings.TrimPrefix(cli.version, "v")
		apiPath = path.Join(cli.basePath, "/v"+v+p)
	} else {
		apiPath = path.Join(cli.basePath, p)
	}

	u := &url.URL{
		Path: apiPath,
	}
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String()
}

// ClientVersion returns the version string associated with this
// instance of the Client. Note that this value can be changed
// via the MALICE_API_VERSION env var.
// This operation doesn't acquire a mutex.
func (cli *Client) ClientVersion() string {
	return cli.version
}

// NegotiateAPIVersion updates the version string associated with this
// instance of the Client to match the latest version the server supports
func (cli *Client) NegotiateAPIVersion(ctx context.Context) {
	ping, _ := cli.Ping(ctx)
	cli.NegotiateAPIVersionPing(ping)
}

// NegotiateAPIVersionPing updates the version string associated with this
// instance of the Client to match the latest version the server supports
func (cli *Client) NegotiateAPIVersionPing(p types.Ping) {
	if cli.manualOverride {
		return
	}

	// try the latest version before versioning headers existed
	if p.APIVersion == "" {
		p.APIVersion = "1.24"
	}

	// if the client is not initialized with a version, start with the latest supported version
	if cli.version == "" {
		cli.version = api.DefaultVersion
	}

	// if server version is lower than the maximum version supported by the Client, downgrade
	if versions.LessThan(p.APIVersion, api.DefaultVersion) {
		cli.version = p.APIVersion
	}
}

// DaemonHost returns the host associated with this instance of the Client.
// This operation doesn't acquire a mutex.
func (cli *Client) DaemonHost() string {
	return cli.host
}

// ParseHost verifies that the given host strings is valid.
func ParseHost(host string) (string, string, string, error) {
	protoAddrParts := strings.SplitN(host, "://", 2)
	if len(protoAddrParts) == 1 {
		return "", "", "", fmt.Errorf("unable to parse docker host `%s`", host)
	}

	var basePath string
	proto, addr := protoAddrParts[0], protoAddrParts[1]
	if proto == "tcp" {
		parsed, err := url.Parse("tcp://" + addr)
		if err != nil {
			return "", "", "", err
		}
		addr = parsed.Host
		basePath = parsed.Path
	}
	return proto, addr, basePath, nil
}

// CustomHTTPHeaders returns the custom http headers associated with this
// instance of the Client. This operation doesn't acquire a mutex.
func (cli *Client) CustomHTTPHeaders() map[string]string {
	m := make(map[string]string)
	for k, v := range cli.customHTTPHeaders {
		m[k] = v
	}
	return m
}

// SetCustomHTTPHeaders updates the custom http headers associated with this
// instance of the Client. This operation doesn't acquire a mutex.
func (cli *Client) SetCustomHTTPHeaders(headers map[string]string) {
	cli.customHTTPHeaders = headers
}
