package plugin

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

// NewPluginCommand returns a cobra command for `plugin` subcommands
// nolint: interfacer
func NewPluginCommand(maliceCli *command.MaliceCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
		Args:  cli.NoArgs,
		RunE:  command.ShowHelp(maliceCli.Err()),
	}

	cmd.AddCommand(
		// newDisableCommand(maliceCli),
		// newEnableCommand(maliceCli),
		// newInspectCommand(maliceCli),
		// newInstallCommand(maliceCli),
		newListCommand(maliceCli),
		// newRemoveCommand(maliceCli),
		// newSetCommand(maliceCli),
		// newPushCommand(maliceCli),
		// newCreateCommand(maliceCli),
		// newUpgradeCommand(maliceCli),
	)
	return cmd
}
