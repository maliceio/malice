package main // import "github.com/maliceio/engine"

import (
	"fmt"
	"os"

	"github.com/maliceio/engine/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
