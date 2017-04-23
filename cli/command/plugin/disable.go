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

type disableOptions struct {
	all bool
}

// NewDisableCommand returns a new cobra disable command for plugins
func NewDisableCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts disableOptions

	cmd := &cobra.Command{
		Use:   "disable [OPTIONS]",
		Short: "Show plugin disable",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runDisable(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin disable")

	return cmd
}

func runDisable(maliceCli *command.MaliceCli, opts disableOptions) (output string, err error) {
	// disableFilters := opts.filter.Value()

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// disable, err := maliceCli.Client().PluginDisable(context.Background(), disableFilters)
	// if err != nil {
	// 	return
	// }

	// if len(disable) > 0 {
	// 	output = "Plugin Disable:\n"
	// 	// TODO: format disable output
	// 	output = disable
	// }

	return
}
