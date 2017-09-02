package watch

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

// NewWatchCommand returns a cobra command for `watch` subcommands
// nolint: interfacer
func NewWatchCommand(maliceCli *command.MaliceCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for files to scan",
		Args:  cli.NoArgs,
		RunE:  command.ShowHelp(maliceCli.Err()),
	}

	cmd.AddCommand()
	return cmd
}
