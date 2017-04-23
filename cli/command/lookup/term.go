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

package lookup

import (
	"fmt"

	"github.com/maliceio/malice/cli/command"
	"github.com/spf13/cobra"
)

type termOptions struct {
	all bool
}

// NewTermCommand returns a new cobra term command for plugins
func NewTermCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts termOptions

	cmd := &cobra.Command{
		Use:   "term [OPTIONS]",
		Short: "Search DB by term",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runTerm(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin term")

	return cmd
}

func runTerm(maliceCli *command.MaliceCli, opts termOptions) (output string, err error) {

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// term, err := maliceCli.Client().LookUpTerm(context.Background(), termFilters)
	// if err != nil {
	// 	return
	// }

	// if len(term) > 0 {
	// 	output = "Plugin Term:\n"
	// 	// TODO: format term output
	// 	output = term
	// }

	return
}
