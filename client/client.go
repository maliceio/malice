package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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

// NewClient initializes a new API client for the given host and API version.
// It uses the given http client as transport.
// It also initializes the custom http headers to add to each request.
func NewClient(host string, version string, client *http.Client, httpHeaders map[string]string) (*Client, error) {
	proto, addr, basePath, err := ParseHost(host)
	if err != nil {
		return nil, err
	}

	if client != nil {
		if _, ok := client.Transport.(*http.Transport); !ok {
			return nil, fmt.Errorf("unable to verify TLS configuration, invalid transport %v", client.Transport)
		}
	} else {
		transport := new(http.Transport)
		client = &http.Client{
			Transport: transport,
		}
	}

	scheme := "http"
	// tlsConfig := resolveTLSConfig(client.Transport)
	// if tlsConfig != nil {
	// 	scheme = "https"
	// }

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
		apiPath = fmt.Sprintf("%s/v%s%s", cli.basePath, v, p)
	} else {
		apiPath = fmt.Sprintf("%s%s", cli.basePath, p)
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
// instance of the Client.
func (cli *Client) ClientVersion() string {
	return cli.version
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
