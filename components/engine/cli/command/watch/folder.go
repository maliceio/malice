package watch

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type folderOpts struct {
	timeout int
	path    string
}

func newWatchFolderCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts folderOpts

	cmd := &cobra.Command{
		Use:   "folder [OPTIONS] INTERFACE",
		Short: "Watch folder for new files",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.path = args[0]
			return runWatchFolder(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runWatchFolder(maliceCli *command.MaliceCli, opts *folderOpts) error {
	return nil
}
