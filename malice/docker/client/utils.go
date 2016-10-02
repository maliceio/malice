package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	er "github.com/maliceio/malice/malice/errors"
)

// Ping pings docker client to see if it is up or not by checking Info.
func (docker *Docker) Ping() bool {

	_, err := docker.Client.Info(context.Background())
	if err != nil {
		er.CheckError(err)
		return false
	}
	return true
}

// parseDockerEndoint returns ip and port from docker endpoint string
func parseDockerEndoint(endpoint string) (string, string, error) {

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return "", "", fmt.Errorf("Unable to parse endpoint: %s", endpoint)
	}

	return hostParts[0], hostParts[1], nil
}
