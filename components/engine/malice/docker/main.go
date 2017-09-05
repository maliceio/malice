package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Client Version: ", cli.ClientVersion())

	responseBody, err := cli.ImagePull(ctx, "malice/elasticsearch", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer responseBody.Close()
	jsonmessage.DisplayJSONMessagesStream(responseBody, os.Stdout, os.Stdout.Fd(), true, nil)

	config := &container.Config{
		Image: "malice/elasticsearch",
	}
	portBindings := nat.PortMap{
		"9200/tcp": {{HostIP: "0.0.0.0", HostPort: "9200"}},
	}
	hostConfig := &container.HostConfig{PortBindings: portBindings}
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, "malice-elasticsearch")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// resultC, errC := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// select {
	// case err := <-errC:
	// 	log.Fatal(err)
	// case result := <-resultC:
	// 	fmt.Println(result.StatusCode)
	// 	// return result.StatusCode, b.String(), nil
	// }

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		// Since       string
		// Timestamps  bool
		Follow: true,
		// Tail        string
	}
	logs, err := cli.ContainerLogs(ctx, resp.ID, options)
	if err != nil {
		panic(err)
	}
	defer logs.Close()

	io.Copy(os.Stdout, logs)
}
