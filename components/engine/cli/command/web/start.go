package web

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type startOpts struct {
	timeout int
	name    string
}

func newStartCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts startOpts

	cmd := &cobra.Command{
		Use:   "start [OPTIONS]",
		Short: "Start web service",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return runStart(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runStart(maliceCli *command.MaliceCli, opts *startOpts) error {
	return nil
}
