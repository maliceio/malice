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
		Short: "Manage watch",
		Args:  cli.NoArgs,
		RunE:  command.ShowHelp(maliceCli.Err()),
	}

	cmd.AddCommand(
		newWatchFolderCommand(maliceCli),
		newNetworkCaptureCommand(maliceCli),
	)
	return cmd
}
