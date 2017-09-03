package plugin

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

func newSetCommand(maliceCli *command.MaliceCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set PLUGIN KEY=VALUE [KEY=VALUE...]",
		Short: "Change settings for a plugin",
		Args:  cli.RequiresMinArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// return maliceCli.Client().PluginSet(context.Background(), args[0], args[1:])
			return nil
		},
	}

	return cmd
}
