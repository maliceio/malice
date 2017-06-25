package client

import (
	"github.com/maliceio/engine/api/types"
	"golang.org/x/net/context"
)

// Ping pings the server and returns the value of the "OS-Type" & "API-Version" headers
func (cli *Client) Ping(ctx context.Context) (types.Ping, error) {
	var ping types.Ping
	req, err := cli.buildRequest("GET", cli.basePath+"/_ping", nil, nil)
	if err != nil {
		return ping, err
	}
	serverResp, err := cli.doRequest(ctx, req)
	if err != nil {
		return ping, err
	}
	defer ensureReaderClosed(serverResp)

	ping.APIVersion = serverResp.header.Get("API-Version")

	ping.OSType = serverResp.header.Get("OSType")

	return ping, nil
}
