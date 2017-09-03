package swarm

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type initOptions struct {
	// swarmOptions
	// listenAddr NodeAddrOption
	// Not a NodeAddrOption because it has no default port.
	advertiseAddr   string
	dataPathAddr    string
	forceNewCluster bool
	availability    string
}

func newInitCommand(maliceCli command.Cli) *cobra.Command {
	opts := initOptions{
	// listenAddr: NewListenAddrOption(),
	}

	cmd := &cobra.Command{
		Use:   "init [OPTIONS]",
		Short: "Initialize a swarm",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(maliceCli, cmd.Flags(), opts)
		},
	}

	flags := cmd.Flags()
	// flags.Var(&opts.listenAddr, flagListenAddr, "Listen address (format: <ip|interface>[:port])")
	// flags.StringVar(&opts.advertiseAddr, flagAdvertiseAddr, "", "Advertised address (format: <ip|interface>[:port])")
	// flags.StringVar(&opts.dataPathAddr, flagDataPathAddr, "", "Address or interface to use for data path traffic (format: <ip|interface>)")
	flags.BoolVar(&opts.forceNewCluster, "force-new-cluster", false, "Force create a new cluster from current state")
	// flags.BoolVar(&opts.autolock, flagAutolock, false, "Enable manager autolocking (requiring an unlock key to start a stopped manager)")
	// flags.StringVar(&opts.availability, flagAvailability, "active", `Availability of the node ("active"|"pause"|"drain")`)
	// addSwarmFlags(flags, &opts.swarmOptions)
	return cmd
}

func runInit(maliceCli command.Cli, flags *pflag.FlagSet, opts initOptions) error {
	return nil
}
