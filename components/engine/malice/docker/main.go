package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Client Version: ", cli.ClientVersion())

	_, err = cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// if _, err = cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning); err != nil {
	// 	panic(err)
	// }

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
