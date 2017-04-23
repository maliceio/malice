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

type hashOptions struct {
	all bool
}

// NewHashCommand returns a new cobra hash command for plugins
func NewHashCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts hashOptions

	cmd := &cobra.Command{
		Use:   "hash [OPTIONS]",
		Short: "Lookup hash (md5/sha1/sha256)",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runHash(maliceCli, opts)
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
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all plugin hash")

	return cmd
}

func runHash(maliceCli *command.MaliceCli, opts hashOptions) (output string, err error) {

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// hash, err := maliceCli.Client().LookUpHash(context.Background(), hashFilters)
	// if err != nil {
	// 	return
	// }

	// if len(hash) > 0 {
	// 	output = "Plugin Hash:\n"
	// 	// TODO: format hash output
	// 	output = hash
	// }

	return
}
