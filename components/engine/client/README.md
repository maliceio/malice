# Go client for the Malice Engine API

The `malice` command uses this package to communicate with the daemon. It can also be used by your own Go applications to do anything the command-line interface does – scanning malware, installing plugins, managing swarms, etc.

For example, to list installed plugins (the equivalent of `malice plugin ls`):

```go
package main

import (
	"context"
	"fmt"

	"github.com/malice/engine/api/types"
	"github.com/malice/engine/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}
```

[Full documentation is available on GoDoc.](https://godoc.org/github.com/malice/engine/client)