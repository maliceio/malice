package plugin

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

type rmOptions struct {
	force bool

	plugins []string
}

func newRemoveCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts rmOptions

	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] PLUGIN [PLUGIN...]",
		Short:   "Remove one or more plugins",
		Aliases: []string{"remove"},
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.plugins = args
			return runRemove(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.force, "force", "f", false, "Force the removal of an active plugin")
	return cmd
}

func runRemove(maliceCli *command.MaliceCli, opts *rmOptions) error {
	// ctx := context.Background()
	//
	// var errs cli.Errors
	// for _, name := range opts.plugins {
	// 	// TODO: pass names to api instead of making multiple api calls
	// 	if err := maliceCli.Client().PluginRemove(ctx, name, types.PluginRemoveOptions{Force: opts.force}); err != nil {
	// 		errs = append(errs, err)
	// 		continue
	// 	}
	// 	fmt.Fprintln(maliceCli.Out(), name)
	// }
	// // Do not simplify to `return errs` because even if errs == nil, it is not a nil-error interface value.
	// if errs != nil {
	// 	return errs
	// }
	return nil
}
