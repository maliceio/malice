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

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/cli"
	"github.com/maliceio/malice/cli/command"
	"github.com/maliceio/malice/cli/command/commands"
	cliconfig "github.com/moby/moby/cli/config"
	cliflags "github.com/moby/moby/cli/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	version string = "dev"
	commit  string = "dev"
	date    string = "dev"
)

func newMaliceCommand(maliceCli *command.MaliceCli) *cobra.Command {
	opts := cliflags.NewClientOptions()
	var flags *pflag.FlagSet

	cmd := &cobra.Command{
		Use:           "malice [OPTIONS] COMMAND [ARG...]",
		Short:         "Open Source Malware Analysis Framework",
		SilenceUsage:  true,
		SilenceErrors: true,
		// TraverseChildren: true,
		// Args:             noArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				showVersion()
				return nil
			}
			return maliceCli.ShowHelp(cmd, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// flags must be the top-level command flags, not cmd.Flags()
			// opts.Common.SetDefaultOptions(flags)
			malicePreRun(opts)
			if err := maliceCli.Initialize(opts); err != nil {
				return err
			}
			return nil
			// return isSupported(cmd, maliceCli.Client().ClientVersion(), maliceCli.OSType(), maliceCli.HasExperimental())
		},
	}
	cli.SetupRootCommand(cmd)

	flags = cmd.Flags()
	flags.BoolVarP(&opts.Version, "version", "v", false, "Print version information and quit")
	// flags.StringVar(&opts.ConfigDir, "config", cliconfig.Dir(), "Location of client config files")
	// opts.Common.InstallFlags(flags)

	// setFlagErrorFunc(maliceCli, cmd, flags, opts)

	// setHelpFunc(maliceCli, cmd, flags, opts)

	cmd.SetOutput(maliceCli.Out())
	// cmd.AddCommand(newDaemonCommand())
	commands.AddCommands(cmd, maliceCli)

	setValidateArgs(maliceCli, cmd, flags, opts)

	return cmd
}

func setValidateArgs(maliceCli *command.MaliceCli, cmd *cobra.Command, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	// The Args is handled by ValidateArgs in cobra, which does not allows a pre-hook.
	// As a result, here we replace the existing Args validation func to a wrapper,
	// where the wrapper will check to see if the feature is supported or not.
	// The Args validation error will only be returned if the feature is supported.
	visitAll(cmd, func(ccmd *cobra.Command) {
		// if there is no tags for a command or any of its parent,
		// there is no need to wrap the Args validation.
		if !hasTags(ccmd) {
			return
		}

		// if ccmd.Args == nil {
		// 	return
		// }

		cmdArgs := cli.NoArgs
		ccmd.Args = func(cmd *cobra.Command, args []string) error {
			initializeMaliceCli(maliceCli, flags, opts)
			return cmdArgs(cmd, args)
		}
	})
}

func initializeMaliceCli(maliceCli *command.MaliceCli, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	if maliceCli.Client() == nil { // when using --help, PersistentPreRun is not called, so initialization is needed.
		// flags must be the top-level command flags, not cmd.Flags()
		// opts.Common.SetDefaultOptions(flags)
		malicePreRun(opts)
		maliceCli.Initialize(opts)
	}
}

// visitAll will traverse all commands from the root.
// This is different from the VisitAll of cobra.Command where only parents
// are checked.
func visitAll(root *cobra.Command, fn func(*cobra.Command)) {
	for _, cmd := range root.Commands() {
		visitAll(cmd, fn)
	}
	fn(root)
}

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf(
		"malice: '%s' is not a malice command.\nSee 'malice --help'", args[0])
}

func main() {
	stdin, stdout, stderr := StdStreams()
	logrus.SetOutput(os.Stderr)

	maliceCli := command.NewMaliceCli(stdin, stdout, stderr)
	cmd := newMaliceCommand(maliceCli)

	if err := cmd.Execute(); err != nil {
		if sterr, ok := err.(cli.StatusError); ok {
			if sterr.Status != "" {
				fmt.Fprintln(os.Stderr, 1)
			}
			// StatusError should only be used for errors, and all errors should
			// have a non-zero exit status, so never exit with 0
			if sterr.StatusCode == 0 {
				os.Exit(1)
			}
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func showVersion() {
	fmt.Printf("Malice version %s, build %s, date %s\n", version, commit, date)
}

func malicePreRun(opts *cliflags.ClientOptions) {
	cliflags.SetLogLevel(opts.Common.LogLevel)

	if opts.ConfigDir != "" {
		cliconfig.SetDir(opts.ConfigDir)
	}

	if opts.Common.Debug {
		// debug.Enable()
	}
}

func isSupported(cmd *cobra.Command) error {
	// clientVersion := details.Client().ClientVersion()
	// osType := details.ServerInfo().OSType

	// errs := []string{}

	// cmd.Flags().VisitAll(func(f *pflag.Flag) {
	// 	if f.Changed {
	// 		if !isVersionSupported(f, clientVersion) {
	// 			errs = append(errs, fmt.Sprintf("\"--%s\" requires API version %s, but the Docker daemon API version is %s", f.Name, getFlagAnnotation(f, "version"), clientVersion))
	// 			return
	// 		}
	// 		if !isOSTypeSupported(f, osType) {
	// 			errs = append(errs, fmt.Sprintf("\"--%s\" requires the Docker daemon to run on %s, but the Docker daemon is running on %s", f.Name, getFlagAnnotation(f, "ostype"), osType))
	// 			return
	// 		}
	// 		if _, ok := f.Annotations["experimental"]; ok && !hasExperimental {
	// 			errs = append(errs, fmt.Sprintf("\"--%s\" is only supported on a Docker daemon with experimental features enabled", f.Name))
	// 		}
	// 	}
	// })
	// if len(errs) > 0 {
	// 	return errors.New(strings.Join(errs, "\n"))
	// }

	return nil
}

// hasTags return true if any of the command's parents has tags
func hasTags(cmd *cobra.Command) bool {
	for curr := cmd; curr != nil; curr = curr.Parent() {
		if len(curr.Tags) > 0 {
			return true
		}
	}

	return false
}

// StdStreams returns the standard streams (stdin, stdout, stderr).
func StdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}
