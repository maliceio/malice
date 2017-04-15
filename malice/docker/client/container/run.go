package container

import (
	"fmt"
	"os"

	"github.com/cloudflare/cfssl/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/cli"
	"github.com/docker/go-connections/nat"

	"github.com/maliceio/malice/malice/docker/client"
	er "github.com/maliceio/malice/malice/errors"
	"golang.org/x/net/context"
)

type runOptions struct {
	autoRemove bool
	detach     bool
	sigProxy   bool
	name       string
	detachKeys string
}

// Run performs a docker run command
func Run(
	docker *client.Docker,
	cmd strslice.StrSlice,
	name string,
	image string,
	logs bool,
	binds []string,
	portBindings nat.PortMap,
	links []string,
	env []string,
) error {
	// stdout, stderr, stdin := os.Stdout, os.Stderr, os.Stdin
	// client := dockerCli.Client()

	// if !opts.detach {
	// 	if err := dockerCli.CheckTtyInput(config.AttachStdin, config.Tty); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	if fl := flags.Lookup("attach"); fl != nil {
	// 		flAttach = fl.Value.(*opttypes.ListOpts)
	// 		if flAttach.Len() != 0 {
	// 			return ErrConflictAttachDetach
	// 		}
	// 	}
	// 	if opts.autoRemove {
	// 		return ErrConflictDetachAutoRemove
	// 	}

	// 	config.AttachStdin = false
	// 	config.AttachStdout = false
	// 	config.AttachStderr = false
	// 	config.StdinOnce = false
	// }

	// Disable sigProxy when in TTY mode
	// if config.Tty {
	// 	opts.sigProxy = false
	// }

	// Telling the Windows daemon the initial size of the tty during start makes
	// a far better user experience rather than relying on subsequent resizes
	// to cause things to catch up.
	// if runtime.GOOS == "windows" {
	// 	hostConfig.ConsoleSize[0], hostConfig.ConsoleSize[1] = dockerCli.GetTtySize()
	// }

	var (
		waitDisplayID chan struct{}
		errCh         chan error
		err           error
	)

	createContConf := &container.Config{
		Image: image,
		Cmd:   cmd,
		Env:   env,
		// Env:   []string{"MALICE_VT_API=" + os.Getenv("MALICE_VT_API")},
	}
	hostConfig := &container.HostConfig{
		// Binds:      []string{maldirs.GetSampledsDir() + ":/malware:ro"},
		// Binds:      []string{"malice:/malware:ro"},
		Binds: binds,
		// NetworkMode:  "malice",
		PortBindings: portBindings,
		Links:        links,
		Privileged:   false,
		AutoRemove:   true,
	}

	networkingConfig := &network.NetworkingConfig{}

	ctx, cancelFun := context.WithCancel(context.Background())

	createResponse, err := createContainer(docker, ctx, createContConf, hostConfig, networkingConfig, hostConfig.ContainerIDFile, name)
	er.CheckError(err)
	// if opts.sigProxy {
	// 	sigc := dockerCli.ForwardAllSignals(ctx, createResponse.ID)
	// 	defer signal.StopCatch(sigc)
	// }

	// if !config.AttachStdout && !config.AttachStderr {
	// Make this asynchronous to allow the client to write to stdin before having to read the ID
	waitDisplayID = make(chan struct{})
	go func() {
		defer close(waitDisplayID)
		fmt.Fprintf(os.Stdout, "%s\n", createResponse.ID)
	}()
	// }
	// if opts.autoRemove && (hostConfig.RestartPolicy.IsAlways() || hostConfig.RestartPolicy.IsOnFailure()) {
	// 	return ErrConflictRestartPolicyAndAutoRemove
	// }
	// attach := config.AttachStdin || config.AttachStdout || config.AttachStderr
	attach := false
	// if attach {
	// 	var (
	// 		out, cerr io.Writer
	// 		in        io.ReadCloser
	// 	)
	// 	if config.AttachStdin {
	// 		in = stdin
	// 	}
	// 	if config.AttachStdout {
	// 		out = stdout
	// 	}
	// 	if config.AttachStderr {
	// 		if config.Tty {
	// 			cerr = stdout
	// 		} else {
	// 			cerr = stderr
	// 		}
	// 	}

	// 	if opts.detachKeys != "" {
	// 		dockerCli.ConfigFile().DetachKeys = opts.detachKeys
	// 	}

	// 	options := types.ContainerAttachOptions{
	// 		Stream:     true,
	// 		Stdin:      config.AttachStdin,
	// 		Stdout:     config.AttachStdout,
	// 		Stderr:     config.AttachStderr,
	// 		DetachKeys: dockerCli.ConfigFile().DetachKeys,
	// 	}

	// 	resp, errAttach := client.ContainerAttach(ctx, createResponse.ID, options)
	// 	if errAttach != nil && errAttach != httputil.ErrPersistEOF {
	// 		// ContainerAttach returns an ErrPersistEOF (connection closed)
	// 		// means server met an error and put it in Hijacked connection
	// 		// keep the error and read detailed error message from hijacked connection later
	// 		return errAttach
	// 	}
	// 	defer resp.Close()

	// 	errCh = promise.Go(func() error {
	// 		errHijack := dockerCli.HoldHijackedConnection(ctx, config.Tty, in, out, cerr, resp)
	// 		if errHijack == nil {
	// 			return errAttach
	// 		}
	// 		return errHijack
	// 	})
	// }

	// if opts.autoRemove {
	// defer func() {
	// 	// Explicitly not sharing the context as it could be "Done" (by calling cancelFun)
	// 	// and thus the container would not be removed.
	// 	if err = Remove(docker, createResponse.ID, true, false, true); err != nil {
	// 		fmt.Fprintf(os.Stderr, "%v\n", err)
	// 	}
	// }()
	// }
	statusChan := waitExitOrRemoved(ctx, docker, createResponse.ID, true)
	//start the container
	log.Debugf("Starting containter: %s", createResponse.ID)
	if err = docker.Client.ContainerStart(ctx, createResponse.ID, types.ContainerStartOptions{}); err != nil {
		// If we have holdHijackedConnection, we should notify
		// holdHijackedConnection we are going to exit and wait
		// to avoid the terminal are not restored.
		if attach {
			cancelFun()
			<-errCh
		}

		<-statusChan

		er.CheckError(err)
	}

	// if (config.AttachStdin || config.AttachStdout || config.AttachStderr) && config.Tty && dockerCli.IsTerminalOut() {
	// 	if err := dockerCli.MonitorTtySize(ctx, createResponse.ID, false); err != nil {
	// 		fmt.Fprintf(stderr, "Error monitoring TTY size: %s\n", err)
	// 	}
	// }

	if errCh != nil {
		if err = <-errCh; err != nil {
			log.Debugf("Error hijack: %s", err)
			return err
		}
	}

	// Detached mode: wait for the id to be displayed and return.
	// if !config.AttachStdout && !config.AttachStderr {
	// Detached mode
	<-waitDisplayID
	// return nil
	// }

	// var status int64

	// Attached mode
	// if opts.autoRemove {
	// Autoremove: wait for the container to finish, retrieve
	// the exit code and remove the container
	// if status, err = docker.Client.ContainerWait(ctx, createResponse.ID); err != nil {
	// 	er.CheckError(err)
	// }
	// if _, status, err = getExitCode(docker, ctx, createResponse.ID); err != nil {
	// 	er.CheckError(err)
	// }
	// } else {
	// 	// No Autoremove: Simply retrieve the exit code
	// 	if !config.Tty && hostConfig.RestartPolicy.IsNone() {
	// 		// In non-TTY mode, we can't detach, so we must wait for container exit
	// 		if status, err = docker.Client.ContainerWait(ctx, createResponse.ID); err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		// In TTY mode, there is a race: if the process dies too slowly, the state could
	// 		// be updated after the getExitCode call and result in the wrong exit code being reported
	// 		if _, status, err = getExitCode(dockerCli, ctx, createResponse.ID); err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	status := <-statusChan
	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}
