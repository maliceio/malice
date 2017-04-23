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

package scan

import (
	"fmt"

	"github.com/maliceio/malice/cli/command"
	"github.com/spf13/cobra"
)

type watchOptions struct {
	all bool
}

// NewWatchCommand returns a new cobra watch command for plugins
func NewWatchCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts watchOptions

	cmd := &cobra.Command{
		Use:   "watch [OPTIONS]",
		Short: "Show plugin watch",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runWatch(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin watch")

	return cmd
}

func runWatch(maliceCli *command.MaliceCli, opts watchOptions) (output string, err error) {

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// watch, err := maliceCli.Client().PluginWatch(context.Background(), watchFilters)
	// if err != nil {
	// 	return
	// }

	// if len(watch) > 0 {
	// 	output = "Plugin Watch:\n"
	// 	// TODO: format watch output
	// 	output = watch
	// }

	return
}
