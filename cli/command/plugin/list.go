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

type listOptions struct {
	all bool
}

// NewListCommand returns a new cobra list command for plugins
func NewListCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts listOptions
	cmd := &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "Show plugin list",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runList(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin list")

	return cmd
}

func runList(maliceCli *command.MaliceCli, opts listOptions) (output string, err error) {
	// listFilters := opts.filter.Value()

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// list, err := maliceCli.Client().PluginList(context.Background(), listFilters)
	// if err != nil {
	// 	return
	// }

	// if len(list) > 0 {
	// 	output = "Plugin List:\n"
	// 	// TODO: format list output
	// 	output = list
	// }

	return
}
