package cli

import (
	"github.com/maliceio/engine/version"

	"github.com/maliceio/engine/command"
	"github.com/maliceio/engine/meta"
	"github.com/mitchellh/cli"
)

// Commands returns the mapping of CLI commands for Vault. The meta
// parameter lets you set meta options for all commands.
func Commands(metaPtr *meta.Meta) map[string]cli.CommandFactory {
	if metaPtr == nil {
		metaPtr = &meta.Meta{
			TokenHelper: command.DefaultTokenHelper,
		}
	}

	// if metaPtr.Ui == nil {
	// 	metaPtr.Ui = &cli.BasicUi{
	// 		Writer:      os.Stdout,
	// 		ErrorWriter: os.Stderr,
	// 	}
	// }

	return map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &command.InitCommand{
				Meta: *metaPtr,
			}, nil
		},

		"ssh": func() (cli.Command, error) {
			return &command.SSHCommand{
				Meta: *metaPtr,
			}, nil
		},

		"version": func() (cli.Command, error) {
			versionInfo := version.GetVersion()

			return &command.VersionCommand{
				VersionInfo: versionInfo,
				Ui:          metaPtr.Ui,
			}, nil
		},
	}
}
