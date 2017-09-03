package search

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type searchOpts struct {
	timeout int
	name    string
}

// NewSearchCommand returns a cobra command for `search` subcommands
// nolint: interfacer
func NewSearchCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts searchOpts

	cmd := &cobra.Command{
		Use:   "search [OPTIONS] HASH",
		Short: "Search for scan results",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return runSearch(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runSearch(maliceCli *command.MaliceCli, opts *searchOpts) error {
	return nil
}
