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

package web

import (
	"github.com/maliceio/malice/cli"
	"github.com/maliceio/malice/cli/command"
	"github.com/spf13/cobra"
)

// NewWebCommand returns a cobra command for `web` subcommands
func NewWebCommand(maliceCli *command.MaliceCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Web services",
		Args:  cli.NoArgs,
		RunE:  maliceCli.ShowHelp,
	}

	cmd.AddCommand(
		NewAPICommand(maliceCli),
		NewUICommand(maliceCli),
	)
	return cmd
}
