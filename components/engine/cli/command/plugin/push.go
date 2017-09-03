package plugin

import (
	"github.com/maliceio/engine/cli"
	"github.com/maliceio/engine/cli/command"
	"github.com/spf13/cobra"
)

func newPushCommand(maliceCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [OPTIONS] PLUGIN[:TAG]",
		Short: "Push a plugin to a registry",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPush(maliceCli, args[0])
		},
	}

	// flags := cmd.Flags()
	// command.AddTrustSigningFlags(flags)

	return cmd
}

func runPush(maliceCli command.Cli, name string) error {
	// named, err := reference.ParseNormalizedNamed(name)
	// if err != nil {
	// 	return err
	// }
	// if _, ok := named.(reference.Canonical); ok {
	// 	return errors.Errorf("invalid name: %s", name)
	// }
	//
	// named = reference.TagNameOnly(named)
	//
	// ctx := context.Background()
	//
	// repoInfo, err := registry.ParseRepositoryInfo(named)
	// if err != nil {
	// 	return err
	// }
	// authConfig := command.ResolveAuthConfig(ctx, maliceCli, repoInfo.Index)
	//
	// encodedAuth, err := command.EncodeAuthToBase64(authConfig)
	// if err != nil {
	// 	return err
	// }
	//
	// responseBody, err := maliceCli.Client().PluginPush(ctx, reference.FamiliarString(named), encodedAuth)
	// if err != nil {
	// 	return err
	// }
	// defer responseBody.Close()
	//
	// if command.IsTrusted() {
	// 	repoInfo.Class = "plugin"
	// 	return image.PushTrustedReference(maliceCli, repoInfo, named, authConfig, responseBody)
	// }
	//
	// return jsonmessage.DisplayJSONMessagesToStream(responseBody, maliceCli.Out(), nil)
	return nil
}
