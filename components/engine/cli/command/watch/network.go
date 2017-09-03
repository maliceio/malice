package watch

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type networkOpts struct {
	timeout int
	name    string
}

func newNetworkCaptureCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts networkOpts

	cmd := &cobra.Command{
		Use:   "network [OPTIONS] INTERFACE",
		Short: "Start network file capture",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return runNetworkCapture(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runNetworkCapture(maliceCli *command.MaliceCli, opts *networkOpts) error {
	return nil
}
