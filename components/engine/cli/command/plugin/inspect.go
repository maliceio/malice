package plugin

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type inspectOptions struct {
	pluginNames []string
	format      string
}

func newInspectCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] PLUGIN [PLUGIN...]",
		Short: "Display detailed information on one or more plugins",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.pluginNames = args
			return runInspect(maliceCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")
	return cmd
}

func runInspect(maliceCli *command.MaliceCli, opts inspectOptions) error {
	// client := maliceCli.Client()
	// ctx := context.Background()
	// getRef := func(ref string) (interface{}, []byte, error) {
	// 	return client.PluginInspectWithRaw(ctx, ref)
	// }
	//
	// return inspect.Inspect(maliceCli.Out(), opts.pluginNames, opts.format, getRef)
	return nil
}
