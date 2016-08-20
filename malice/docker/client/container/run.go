package container

import (
	"fmt"
	"os"

	"github.com/cloudflare/cfssl/log"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	networktypes "github.com/docker/engine-api/types/network"
	"github.com/maliceio/malice/malice/docker/client"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/spf13/pflag"
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
func Run(docker *client.Docker, flags *pflag.FlagSet, opts *runOptions, copts *runconfigopts.ContainerOptions) error {
	// stdout, stderr, stdin := os.Stdout, os.Stderr, os.Stdin
	stderr := os.Stderr
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
	config := container.Config{
	// Hostname        string                // Hostname
	// Domainname      string                // Domainname
	// User            string                // User that will run the command(s) inside the container, also support user:group
	// AttachStdin     bool                  // Attach the standard input, makes possible user interaction
	// AttachStdout    bool                  // Attach the standard output
	// AttachStderr    bool                  // Attach the standard error
	// ExposedPorts    map[nat.Port]struct{} `json:",omitempty"` // List of exposed ports
	// Tty             bool                  // Attach standard streams to a tty, including stdin if it is not closed.
	// OpenStdin       bool                  // Open stdin
	// StdinOnce       bool                  // If true, close stdin after the 1 attached client disconnects.
	// Env             []string              // List of environment variable to set in the container
	// Cmd             strslice.StrSlice     // Command to run when starting the container
	// Healthcheck     *HealthConfig         `json:",omitempty"` // Healthcheck describes how to check the container is healthy
	// ArgsEscaped     bool                  `json:",omitempty"` // True if command is already escaped (Windows specific)
	// Image           string                // Name of the image as it was passed by the operator (eg. could be symbolic)
	// Volumes         map[string]struct{}   // List of volumes (mounts) used for the container
	// WorkingDir      string                // Current directory (PWD) in the command will be launched
	// Entrypoint      strslice.StrSlice     // Entrypoint to run when starting the container
	// NetworkDisabled bool                  `json:",omitempty"` // Is network disabled
	// MacAddress      string                `json:",omitempty"` // Mac Address of the container
	// OnBuild         []string              // ONBUILD metadata that were defined on the image Dockerfile
	// Labels          map[string]string     // List of labels set to this container
	// StopSignal      string                `json:",omitempty"` // Signal to stop a container
	// StopTimeout     *int                  `json:",omitempty"` // Timeout (in seconds) to stop a container
	// Shell           strslice.StrSlice     `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT
	}
	hostConfig := container.HostConfig{
	// // Applicable to all platforms
	// Binds           []string      // List of volume bindings for this container
	// ContainerIDFile string        // File (path) where the containerId is written
	// LogConfig       LogConfig     // Configuration of the logs for this container
	// NetworkMode     NetworkMode   // Network mode to use for the container
	// PortBindings    nat.PortMap   // Port mapping between the exposed port (container) and the host
	// RestartPolicy   RestartPolicy // Restart policy to be used for the container
	// AutoRemove      bool          // Automatically remove container when it exits
	// VolumeDriver    string        // Name of the volume driver used to mount volumes
	// VolumesFrom     []string      // List of volumes to take from other container
	//
	// // Applicable to UNIX platforms
	// CapAdd          strslice.StrSlice // List of kernel capabilities to add to the container
	// CapDrop         strslice.StrSlice // List of kernel capabilities to remove from the container
	// DNS             []string          `json:"Dns"`        // List of DNS server to lookup
	// DNSOptions      []string          `json:"DnsOptions"` // List of DNSOption to look for
	// DNSSearch       []string          `json:"DnsSearch"`  // List of DNSSearch to look for
	// ExtraHosts      []string          // List of extra hosts
	// GroupAdd        []string          // List of additional groups that the container process will run as
	// IpcMode         IpcMode           // IPC namespace to use for the container
	// Cgroup          CgroupSpec        // Cgroup to use for the container
	// Links           []string          // List of links (in the name:alias form)
	// OomScoreAdj     int               // Container preference for OOM-killing
	// PidMode         PidMode           // PID namespace to use for the container
	// Privileged      bool              // Is the container in privileged mode
	// PublishAllPorts bool              // Should docker publish all exposed port for the container
	// ReadonlyRootfs  bool              // Is the container root filesystem in read-only
	// SecurityOpt     []string          // List of string values to customize labels for MLS systems, such as SELinux.
	// StorageOpt      map[string]string `json:",omitempty"` // Storage driver options per container.
	// Tmpfs           map[string]string `json:",omitempty"` // List of tmpfs (mounts) used for the container
	// UTSMode         UTSMode           // UTS namespace to use for the container
	// UsernsMode      UsernsMode        // The user namespace to use for the container
	// ShmSize         int64             // Total shm memory usage
	// Sysctls         map[string]string `json:",omitempty"` // List of Namespaced sysctls used for the container
	// Runtime         string            `json:",omitempty"` // Runtime to use with this container
	//
	// // Applicable to Windows
	// ConsoleSize [2]int    // Initial console size
	// Isolation   Isolation // Isolation technology of the container (eg default, hyperv)
	//
	// // Contains container's resources (cgroups, ulimits)
	// Resources
	//
	// // Mounts specs used by the container
	// Mounts []mount.Mount `json:",omitempty"`
	}
	networkingConfig := networktypes.NetworkingConfig{
	// EndpointsConfig map[string]*EndpointSettings // Endpoint configs for each connecting network
	}
	ctx, cancelFun := context.WithCancel(context.Background())

	createResponse, err := createContainer(docker, ctx, &config, &hostConfig, &networkingConfig, hostConfig.ContainerIDFile, opts.name)
	er.CheckError(err)
	// if opts.sigProxy {
	// 	sigc := dockerCli.ForwardAllSignals(ctx, createResponse.ID)
	// 	defer signal.StopCatch(sigc)
	// }
	var (
		waitDisplayID chan struct{}
		errCh         chan error
	)
	// if !config.AttachStdout && !config.AttachStderr {
	// 	// Make this asynchronous to allow the client to write to stdin before having to read the ID
	// 	waitDisplayID = make(chan struct{})
	// 	go func() {
	// 		defer close(waitDisplayID)
	// 		fmt.Fprintf(stdout, "%s\n", createResponse.ID)
	// 	}()
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
	defer func() {
		// Explicitly not sharing the context as it could be "Done" (by calling cancelFun)
		// and thus the container would not be removed.
		if err := Remove(docker, createResponse.ID, true, false, true); err != nil {
			fmt.Fprintf(stderr, "%v\n", err)
		}
	}()
	// }

	//start the container
	if err := docker.Client.ContainerStart(ctx, createResponse.ID, types.ContainerStartOptions{}); err != nil {
		// If we have holdHijackedConnection, we should notify
		// holdHijackedConnection we are going to exit and wait
		// to avoid the terminal are not restored.
		if attach {
			cancelFun()
			<-errCh
		}

		er.CheckError(err)
	}

	// if (config.AttachStdin || config.AttachStdout || config.AttachStderr) && config.Tty && dockerCli.IsTerminalOut() {
	// 	if err := dockerCli.MonitorTtySize(ctx, createResponse.ID, false); err != nil {
	// 		fmt.Fprintf(stderr, "Error monitoring TTY size: %s\n", err)
	// 	}
	// }

	if errCh != nil {
		if err := <-errCh; err != nil {
			log.Debugf("Error hijack: %s", err)
			return err
		}
	}

	// Detached mode: wait for the id to be displayed and return.
	if !config.AttachStdout && !config.AttachStderr {
		// Detached mode
		<-waitDisplayID
		return nil
	}

	var status int

	// Attached mode
	// if opts.autoRemove {
	// Autoremove: wait for the container to finish, retrieve
	// the exit code and remove the container
	if status, err = docker.Client.ContainerWait(ctx, createResponse.ID); err != nil {
		er.CheckError(err)
	}
	if _, status, err = getExitCode(docker, ctx, createResponse.ID); err != nil {
		er.CheckError(err)
	}
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
	if status != 0 {
		return fmt.Errorf("Status: %d", status)
	}
	return nil
}
