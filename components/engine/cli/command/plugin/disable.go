package plugin

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

func newDisableCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "disable [OPTIONS] PLUGIN",
		Short: "Disable a plugin",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDisable(maliceCli, args[0], force)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Force the disable of an active plugin")
	return cmd
}

func runDisable(maliceCli *command.MaliceCli, name string, force bool) error {
	// if err := maliceCli.Client().PluginDisable(context.Background(), name, types.PluginDisableOptions{Force: force}); err != nil {
	// 	return err
	// }
	// fmt.Fprintln(maliceCli.Out(), name)
	return nil
}
