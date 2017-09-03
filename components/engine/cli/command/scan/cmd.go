package scan

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type scanOpts struct {
	timeout int
	name    string
}

// NewScanCommand returns a cobra command for `search` subcommands
// nolint: interfacer
func NewScanCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts scanOpts

	cmd := &cobra.Command{
		Use:   "scan [OPTIONS] PLUGIN",
		Short: "Scan a file",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return runScan(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runScan(maliceCli *command.MaliceCli, opts *scanOpts) error {
	return nil
}
