package system

import (
	"github.com/docker/docker/pkg/templates"
	"github.com/maliceio/malice/cli"
	"github.com/maliceio/malice/cli/command"
	"github.com/spf13/cobra"
)

var versionTemplate = `Client:
 Version:      {{.Client.Version}}
 API version:  {{.Client.APIVersion}}
 Go version:   {{.Client.GoVersion}}
 Git commit:   {{.Client.GitCommit}}
 Built:        {{.Client.BuildTime}}
 OS/Arch:      {{.Client.Os}}/{{.Client.Arch}}{{if .ServerOK}}
Server:
 Version:      {{.Server.Version}}
 API version:  {{.Server.APIVersion}} (minimum version {{.Server.MinAPIVersion}})
 Go version:   {{.Server.GoVersion}}
 Git commit:   {{.Server.GitCommit}}
 Built:        {{.Server.BuildTime}}
 OS/Arch:      {{.Server.Os}}/{{.Server.Arch}}
 Experimental: {{.Server.Experimental}}{{end}}`

type versionOptions struct {
	format string
}

// NewVersionCommand creates a new cobra.Command for `docker version`
func NewVersionCommand(maliceCli *command.MaliceCli) *cobra.Command {
	var opts versionOptions

	cmd := &cobra.Command{
		Use:   "version [OPTIONS]",
		Short: "Show the Malice version information",
		// Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(maliceCli, &opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")

	return cmd
}

func runVersion(maliceCli *command.MaliceCli, opts *versionOptions) error {
	// ctx := context.Background()

	templateFormat := versionTemplate
	if opts.format != "" {
		templateFormat = opts.format
	}

	tmpl, err := templates.Parse(templateFormat)
	if err != nil {
		return cli.StatusError{StatusCode: 64,
			Status: "Template parsing error: " + err.Error()}
	}

	// APIVersion := maliceCli.Client().ClientVersion()
	// if defaultAPIVersion := dockerCli.DefaultVersion(); APIVersion != defaultAPIVersion {
	// 	APIVersion = fmt.Sprintf("%s (downgraded from %s)", APIVersion, defaultAPIVersion)
	// }

	// vd := types.VersionResponse{
	// 	Client: &types.Version{
	// 		Version:    dockerversion.Version,
	// 		APIVersion: APIVersion,
	// 		GoVersion:  runtime.Version(),
	// 		GitCommit:  dockerversion.GitCommit,
	// 		BuildTime:  dockerversion.BuildTime,
	// 		Os:         runtime.GOOS,
	// 		Arch:       runtime.GOARCH,
	// 	},
	// }

	// serverVersion, err := maliceCli.Client().ServerVersion(ctx)
	// if err == nil {
	// 	vd.Server = &serverVersion
	// }

	// // first we need to make BuildTime more human friendly
	// t, errTime := time.Parse(time.RFC3339Nano, vd.Client.BuildTime)
	// if errTime == nil {
	// 	vd.Client.BuildTime = t.Format(time.ANSIC)
	// }

	// if vd.ServerOK() {
	// 	t, errTime = time.Parse(time.RFC3339Nano, vd.Server.BuildTime)
	// 	if errTime == nil {
	// 		vd.Server.BuildTime = t.Format(time.ANSIC)
	// 	}
	// }

	// if err2 := tmpl.Execute(maliceCli.Out(), vd); err2 != nil && err == nil {
	if err2 := tmpl.Execute(maliceCli.Out(), nil); err2 != nil && err == nil {
		err = err2
	}
	maliceCli.Out().Write([]byte{'\n'})
	return err
}
