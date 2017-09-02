package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/maliceio/engine/cli/command/commands"
	cliconfig "github.com/maliceio/engine/cli/config"
	cliflags "github.com/maliceio/engine/cli/flags"
	"github.com/maliceio/engine/malice/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newMaliceCommand(maliceCLI *command.MaliceCli) *cobra.Command {
	opts := cliflags.NewClientOptions()
	var flags *pflag.FlagSet

	cmd := &cobra.Command{
		Use:           "malice [OPTIONS] COMMAND [ARG...]",
		Short:         "Open Source Malware Analysis Framework",
		SilenceUsage:  true,
		SilenceErrors: true,
		// TraverseChildren: true,
		Args: noArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				showVersion()
				return nil
			}
			return command.ShowHelp(maliceCLI.Err())(cmd, args)
		},
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	// daemon command is special, we redirect directly to another binary
		// 	if cmd.Name() == "daemon" {
		// 		return nil
		// 	}
		// 	// flags must be the top-level command flags, not cmd.Flags()
		// 	opts.Common.SetDefaultOptions(flags)
		// 	malicePreRun(opts)
		// 	if err := maliceCLI.Initialize(opts); err != nil {
		// 		return err
		// 	}
		// 	return nil
		// },
	}
	cli.SetupRootCommand(cmd)

	flags = cmd.Flags()
	flags.BoolVarP(&opts.Version, "version", "v", false, "Print version information and quit")
	flags.StringVar(&opts.ConfigDir, "config", cliconfig.Dir(), "Location of client config files")
	opts.Common.InstallFlags(flags)

	// setFlagErrorFunc(maliceCLI, cmd, flags, opts)

	// setHelpFunc(maliceCLI, cmd, flags, opts)

	// cmd.SetOutput(maliceCLI.Out())
	// cmd.AddCommand(newDaemonCommand())
	commands.AddCommands(cmd, maliceCLI)

	// setValidateArgs(maliceCLI, cmd, flags, opts)
	return cmd
}

func setHelpFunc(maliceCli *command.MaliceCli, cmd *cobra.Command, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	cmd.SetHelpFunc(func(ccmd *cobra.Command, args []string) {
		initializeMaliceCli(maliceCli, flags, opts)

		if err := ccmd.Help(); err != nil {
			ccmd.Println(err)
		}
	})
}

func initializeMaliceCli(maliceCLI *command.MaliceCli, flags *pflag.FlagSet, opts *cliflags.ClientOptions) {
	if maliceCLI.Client() == nil { // when using --help, PersistentPreRun is not called, so initialization is needed.
		// flags must be the top-level command flags, not cmd.Flags()
		opts.Common.SetDefaultOptions(flags)
		malicePreRun(opts)
		maliceCLI.Initialize(opts)
	}
}

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf(
		"malice: '%s' is not a malice command.\nSee 'malice --help'", args[0])
}

func main() {
	logrus.SetOutput(os.Stderr)

	maliceCLI := command.NewMaliceCli(os.Stderr)
	cmd := newMaliceCommand(maliceCLI)

	if err := cmd.Execute(); err != nil {
		if sterr, ok := err.(cli.StatusError); ok {
			if sterr.Status != "" {
				fmt.Fprintln(os.Stderr, sterr.Status)
			}
			// StatusError should only be used for errors, and all errors should
			// have a non-zero exit status, so never exit with 0
			if sterr.StatusCode == 0 {
				os.Exit(1)
			}
			os.Exit(sterr.StatusCode)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func showVersion() {
	fmt.Printf("Malice version %s, build %s\n", version.Version, version.GitCommit)
}

func malicePreRun(opts *cliflags.ClientOptions) {
	cliflags.SetLogLevel(opts.Common.LogLevel)

	if opts.ConfigDir != "" {
		cliconfig.SetDir(opts.ConfigDir)
	}

	// if opts.Common.Debug {
	// 	debug.Enable()
	// }
}

// // hasTags return true if any of the command's parents has tags
// func hasTags(cmd *cobra.Command) bool {
// 	for curr := cmd; curr != nil; curr = curr.Parent() {
// 		if len(curr.Tags) > 0 {
// 			return true
// 		}
// 	}

// 	return false
// }
