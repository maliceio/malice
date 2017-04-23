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

type enableOptions struct {
	all bool
}

// NewEnableCommand returns a new cobra enable command for plugins
func NewEnableCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts enableOptions

	cmd := &cobra.Command{
		Use:   "enable [OPTIONS]",
		Short: "Show plugin enable",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runEnable(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin enable")

	return cmd
}

func runEnable(maliceCli *command.MaliceCli, opts enableOptions) (output string, err error) {
	// enableFilters := opts.filter.Value()

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// enable, err := maliceCli.Client().PluginEnable(context.Background(), enableFilters)
	// if err != nil {
	// 	return
	// }

	// if len(enable) > 0 {
	// 	output = "Plugin Enable:\n"
	// 	// TODO: format enable output
	// 	output = enable
	// }

	return
}
