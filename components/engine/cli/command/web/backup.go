package web

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type backupOpts struct {
	timeout int
	name    string
}

func newBackUpCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts backupOpts

	cmd := &cobra.Command{
		Use:   "backup [OPTIONS] FILE",
		Short: "Backup malice database",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return runBackUp(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runBackUp(maliceCli *command.MaliceCli, opts *backupOpts) error {
	return nil
}
