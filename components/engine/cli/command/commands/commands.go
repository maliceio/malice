package commands

import (
	"github.com/maliceio/engine/cli/command/plugin"
	"github.com/spf13/cobra"
)

// AddCommands adds all the commands from cli/command to the root command
func AddCommands(cmd *cobra.Command, dockerCli *command.DockerCli) {
	cmd.AddCommand(
		// plugin
		plugin.NewPluginCommand(dockerCli),
	)

}
