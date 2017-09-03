package commands

import (
	"github.com/maliceio/engine/cli/command"
	"github.com/maliceio/engine/cli/command/plugin"
	"github.com/maliceio/engine/cli/command/scan"
	"github.com/spf13/cobra"
)

// AddCommands adds all the commands from cli/command to the root command
func AddCommands(cmd *cobra.Command, maliceCli *command.MaliceCli) {
	cmd.AddCommand(
		// plugin
		plugin.NewPluginCommand(maliceCli),

		// scan
		scan.NewScanCommand(maliceCli),
	)

}
