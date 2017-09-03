package swarm

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

// NewSwarmCommand returns a cobra command for `swarm` subcommands
// nolint: interfacer
func NewSwarmCommand(maliceCli *command.MaliceCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swarm",
		Short: "Manage swarm",
		Args:  cli.NoArgs,
		RunE:  command.ShowHelp(maliceCli.Err()),
	}

	cmd.AddCommand(
		newInitCommand(maliceCli),
	)
	return cmd
}
