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

type installOptions struct {
	all bool
}

// NewInstallCommand returns a new cobra install command for plugins
func NewInstallCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts installOptions
	cmd := &cobra.Command{
		Use:   "install [OPTIONS]",
		Short: "Show plugin install",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runInstall(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin install")

	return cmd
}

func runInstall(maliceCli *command.MaliceCli, opts installOptions) (output string, err error) {
	// installFilters := opts.filter.Value()

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// install, err := maliceCli.Client().PluginInstall(context.Background(), installFilters)
	// if err != nil {
	// 	return
	// }

	// if len(install) > 0 {
	// 	output = "Plugin Install:\n"
	// 	// TODO: format install output
	// 	output = install
	// }

	return
}
