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

type updateOptions struct {
	all bool
}

// NewUpdateCommand returns a new cobra update command for plugins
func NewUpdateCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts updateOptions
	cmd := &cobra.Command{
		Use:   "update [OPTIONS]",
		Short: "Show plugin update",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runUpdate(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin update")

	return cmd
}

func runUpdate(maliceCli *command.MaliceCli, opts updateOptions) (output string, err error) {
	// updateFilters := opts.filter.Value()

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// update, err := maliceCli.Client().PluginUpdate(context.Background(), updateFilters)
	// if err != nil {
	// 	return
	// }

	// if len(update) > 0 {
	// 	output = "Plugin Update:\n"
	// 	// TODO: format update output
	// 	output = update
	// }

	return
}
