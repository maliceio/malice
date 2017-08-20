// Copyright © 2017 blacktop <https://github.com/maliceio>
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

package cmd

import (
	"fmt"
	"os"

	"github.com/maliceio/engine/daemon/config"
	"github.com/maliceio/engine/malice/version"
	"github.com/spf13/cobra"
)

const (
	// DefaultCaFile is the default filename for the CA pem file
	DefaultCaFile = "ca.pem"
	// DefaultKeyFile is the default filename for the key pem file
	DefaultKeyFile = "key.pem"
	// DefaultCertFile is the default filename for the cert pem file
	DefaultCertFile = "cert.pem"
)

var (
	maliceCertPath  = os.Getenv("MALICE_CERT_PATH")
	maliceTLSVerify = os.Getenv("MALICE_TLS_VERIFY") != ""
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:           "daemon",
	Short:         "Start the malice daemon",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := newDaemonOptions(config.New())
		installConfigFlags(opts.daemonConfig, cmd.Flags())

		return runDaemon(opts)
	},
}

func init() {
	RootCmd.AddCommand(daemonCmd)

	daemonCmd.Flags().BoolVarP("version", "v", false, "Show malice version")
}

func runDaemon(opts) error {
	if daemonCmd.Flags().GetBool("version") {
		showVersion()
		return nil
	}

	daemon := NewDaemon()

	err = daemon.start(opts)
	notifyShutdown(err)
	return err
}

func showVersion() {
	fmt.Printf("Malice version %s, build %s\n", version.Version, version.GitCommit)
}
