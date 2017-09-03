package web

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type stopOpts struct {
	timeout int
	name    string
}

func newStopCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts stopOpts

	cmd := &cobra.Command{
		Use:   "stop [OPTIONS]",
		Short: "Stop web service",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return runStop(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runStop(maliceCli *command.MaliceCli, opts *stopOpts) error {
	return nil
}
