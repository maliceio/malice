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

type infoOptions struct {
	all bool
}

// NewInfoCommand returns a new cobra info command for plugins
func NewInfoCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts infoOptions

	cmd := &cobra.Command{
		Use:   "info [OPTIONS]",
		Short: "Show plugin info",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runInfo(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin info")

	return cmd
}

func runInfo(maliceCli *command.MaliceCli, opts infoOptions) (output string, err error) {

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// info, err := maliceCli.Client().PluginInfo(context.Background(), infoFilters)
	// if err != nil {
	// 	return
	// }

	// if len(info) > 0 {
	// 	output = "Plugin Info:\n"
	// 	// TODO: format info output
	// 	output = info
	// }

	return
}
