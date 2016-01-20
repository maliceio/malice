package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
)

func main() {
	// defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	// cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.21", nil, defaultHeaders)
	// os.Unsetenv("DOCKER_TLS_VERIFY")
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(options)
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		fmt.Println(c.ID)
	}

	pullOptions := types.ImagePullOptions{
		ImageID: "malice/fprot",
		Tag:     "latest",
		// RegistryAuth: encodedAuth,
	}

	responseBody, err := cli.ImagePull(pullOptions, nil)
	if err != nil {
		panic(err)
	}
	defer responseBody.Close()

	jsonmessage.DisplayJSONMessagesStream(responseBody, os.Stdout, os.Stdout.Fd(), true, nil)
}
