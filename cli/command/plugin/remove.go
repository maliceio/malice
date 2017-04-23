// Copyright © 2017 blacktop <https://github.com/blacktop>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

import (
	"fmt"

	"github.com/maliceio/malice/cli/command"
	"github.com/spf13/cobra"
)

type removeOptions struct {
	all bool
}

// NewRemoveCommand returns a new cobra remove command for plugins
func NewRemoveCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts removeOptions
	cmd := &cobra.Command{
		Use:   "remove [OPTIONS]",
		Short: "Show plugin remove",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runRemove(maliceCli, opts)
			if err != nil {
				return err
			}
			if output != "" {
				fmt.Fprintln(maliceCli.Out(), output)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin remove")

	return cmd
}

func runRemove(maliceCli *command.MaliceCli, opts removeOptions) (output string, err error) {
	// removeFilters := opts.filter.Value()

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// remove, err := maliceCli.Client().PluginRemove(context.Background(), removeFilters)
	// if err != nil {
	// 	return
	// }

	// if len(remove) > 0 {
	// 	output = "Plugin Remove:\n"
	// 	// TODO: format remove output
	// 	output = remove
	// }

	return
}
