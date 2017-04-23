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

type scanOptions struct {
	force bool
}

// NewScanCommand returns a cobra command for `plugin` subcommands
func NewScanCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts scanOptions

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan a file",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runScan(maliceCli, opts)
			if err != nil {
				return err
			}
			if output != "" {
				fmt.Fprintln(maliceCli.Out(), output)
			}

			return nil
		},
	}

	cmd.AddCommand(
		NewWatchCommand(maliceCli),
	)

	flags := cmd.Flags()
	flags.BoolVarP(&opts.force, "force", "f", false, "Rescan file")

	return cmd
}

func runScan(maliceCli *command.MaliceCli, opts scanOptions) (output string, err error) {

	// if !opts.force && !command.PromptForConfirmation(maliceCli.In(), maliceCli.Out(), warning) {
	// 	return
	// }

	// scan, err := maliceCli.Client().Scan(context.Background(), opts.force)
	// if err != nil {
	// 	return
	// }

	// if len(scan) > 0 {
	// 	output = "Scan:\n"
	// 	// TODO: format scan output
	// 	output = scan
	// }

	return
}
