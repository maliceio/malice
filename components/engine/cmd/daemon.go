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
	"path/filepath"

	"github.com/maliceio/engine/malice/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgDaemonFile = os.Getenv("MALICE_CONFIG")
	// Version version output flag
	Version bool
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:           "daemon",
	Short:         "Start the malice daemon",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// opts := newDaemonOptions(config.New())
		// installConfigFlags(opts.daemonConfig, cmd.Flags())

		return runDaemon()
	},
}

func init() {
	initDaemonConfig()
	RootCmd.AddCommand(daemonCmd)

	daemonCmd.Flags().BoolVarP(&Version, "version", "v", false, "Show malice version")
}

func runDaemon() error {
	if Version {
		showVersion()
		return nil
	}

	// daemon := NewDaemon()

	// err = daemon.start(opts)
	// notifyShutdown(err)
	// return err
	return nil
}

func showVersion() {
	fmt.Printf("Malice version %s, build %s\n", version.Version, version.GitCommit)
}

// initConfig reads in config file and ENV variables if set.
func initDaemonConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("toml")
	viper.SetConfigName(filepath.Join(cfgDaemonFile, "config", "daemon.toml")) // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("HOME"))                                     // adding home directory as first search path
	viper.AutomaticEnv()                                                       // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
