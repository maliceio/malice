Go client for the Malice Engine API
===================================

The `malice` command uses this package to communicate with the daemon. It can also be used by your own Go applications to do anything the command-line interface does – running scans, installing plugins, managing swarms, etc.

For example, to list installed plugins (the equivalent of `malice plugin ls`\):

```go
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
```

[Full documentation is available on GoDoc.](https://godoc.org/github.com/maliceio/engine/client)
