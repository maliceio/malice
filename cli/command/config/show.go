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

package config

import (
	"fmt"

	"github.com/maliceio/malice/cli/command"
	"github.com/spf13/cobra"
)

type showOptions struct {
	all bool
}

// NewShowCommand returns a new cobra show command for plugins
func NewShowCommand(maliceCli *command.MaliceCli) *cobra.Command {

	var opts showOptions

	cmd := &cobra.Command{
		Use:   "show [OPTIONS]",
		Short: "Show plugin show",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runShow(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin show")

	return cmd
}

func runShow(maliceCli *command.MaliceCli, opts showOptions) (output string, err error) {

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// show, err := maliceCli.Client().ConfigShow(context.Background(), opts)
	// if err != nil {
	// 	return
	// }

	// if len(show) > 0 {
	// 	output = "Plugin Show:\n"
	// 	// TODO: format show output
	// 	output = show
	// }

	return
}
