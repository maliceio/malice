package plugin

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type enableOpts struct {
	timeout int
	name    string
}

func newEnableCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts enableOpts

	cmd := &cobra.Command{
		Use:   "enable [OPTIONS] PLUGIN",
		Short: "Enable a plugin",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return runEnable(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.timeout, "timeout", 0, "HTTP client timeout (in seconds)")
	return cmd
}

func runEnable(maliceCli *command.MaliceCli, opts *enableOpts) error {
	// name := opts.name
	// if opts.timeout < 0 {
	// 	return errors.Errorf("negative timeout %d is invalid", opts.timeout)
	// }
	//
	// if err := maliceCli.Client().PluginEnable(context.Background(), name, types.PluginEnableOptions{Timeout: opts.timeout}); err != nil {
	// 	return err
	// }
	// fmt.Fprintln(maliceCli.Out(), name)
	return nil
}
