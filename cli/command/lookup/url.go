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

type urlOptions struct {
	all bool
}

// NewURLCommand returns a new cobra url command for plugins
func NewURLCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts urlOptions

	cmd := &cobra.Command{
		Use:   "url [OPTIONS]",
		Short: "Lookup URL/Domain",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runURL(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin url")

	return cmd
}

func runURL(maliceCli *command.MaliceCli, opts urlOptions) (output string, err error) {

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// url, err := maliceCli.Client().LookUpURL(context.Background(), urlFilters)
	// if err != nil {
	// 	return
	// }

	// if len(url) > 0 {
	// 	output = "LookUp URL:\n"
	// 	// TODO: format url output
	// 	output = url
	// }

	return
}
