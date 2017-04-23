package commands

import (
	"github.com/maliceio/malice/cli/command"
	"github.com/maliceio/malice/cli/command/config"
	"github.com/maliceio/malice/cli/command/lookup"
	"github.com/maliceio/malice/cli/command/plugin"
	"github.com/maliceio/malice/cli/command/scan"
	"github.com/maliceio/malice/cli/command/system"
	"github.com/spf13/cobra"
)

// AddCommands adds all the commands from cli/command to the root command
func AddCommands(cmd *cobra.Command, maliceCli *command.MaliceCli) {
	cmd.AddCommand(
		// config
		config.NewConfigCommand(maliceCli),
		// scan
		scan.NewScanCommand(maliceCli),
		// lookup
		lookup.NewLookUpCommand(maliceCli),
		// plugin
		plugin.NewPluginCommand(maliceCli),
		// system
		system.NewVersionCommand(maliceCli),
	)
}
