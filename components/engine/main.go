package main // import "github.com/maliceio/engine"

import (
	"os"

	"github.com/maliceio/engine/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
