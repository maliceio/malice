package container

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	cont "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	apiclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/registry"
	"github.com/maliceio/malice/malice/docker/client"
	er "github.com/maliceio/malice/malice/errors"
	"golang.org/x/net/context"
)

func pullImage(ctx context.Context, docker *client.Docker, image string, out io.Writer) error {
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return err
	}

	// Resolve the Repository name from fqn to RepositoryInfo
	_, err = registry.ParseRepositoryInfo(ref)
	if err != nil {
		return err
	}

	// authConfig := dockerCli.ResolveAuthConfig(ctx, repoInfo.Index)
	// encodedAuth, err := client.EncodeAuthToBase64(authConfig)
	// if err != nil {
	// 	return err
	// }

	options := types.ImageCreateOptions{
	// RegistryAuth: encodedAuth,
	}

	responseBody, err := docker.Client.ImageCreate(ctx, image, options)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	return jsonmessage.DisplayJSONMessagesStream(
		responseBody,
		out,
		os.Stdout.Fd(),
		true,
		nil)
}

type cidFile struct {
	path    string
	file    *os.File
	written bool
}

func (cid *cidFile) Close() error {
	cid.file.Close()

	if !cid.written {
		if err := os.Remove(cid.path); err != nil {
			return fmt.Errorf("failed to remove the CID file '%s': %s \n", cid.path, err)
		}
	}

	return nil
}

func (cid *cidFile) Write(id string) error {
	if _, err := cid.file.Write([]byte(id)); err != nil {
		return fmt.Errorf("Failed to write the container ID to the file: %s", err)
	}
	cid.written = true
	return nil
}

func newCIDFile(path string) (*cidFile, error) {
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("Container ID file found, make sure the other container isn't running or delete %s", path)
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the container ID file: %s", err)
	}

	return &cidFile{path: path, file: f}, nil
}

// createContainer
func createContainer(docker *client.Docker, ctx context.Context, config *cont.Config, hostConfig *cont.HostConfig, networkingConfig *networktypes.NetworkingConfig, cidfile, name string) (cont.ContainerCreateCreatedBody, error) {
	stderr := os.Stderr
	// log.Info("cidfile: ", cidfile)
	var containerIDFile *cidFile
	// log.Info("containerIDFile: ", containerIDFile)
	if cidfile != "" {
		var err error
		if containerIDFile, err = newCIDFile(cidfile); err != nil {
			er.CheckError(err)
			return cont.ContainerCreateCreatedBody{}, err
		}
		// log.Info("NEW containerIDFile: ", containerIDFile)
		defer containerIDFile.Close()
	}

	// var trustedRef reference.Canonical
	// _, ref, err := reference.ParseIDOrReference(config.Image)
	// if err != nil {
	// 	return nil, err
	// }
	// if ref != nil {
	// 	ref = reference.WithDefaultTag(ref)

	// 	if ref, ok := ref.(reference.NamedTagged); ok {
	// 		var err error
	// 		trustedRef, err = docker.TrustedReference(ctx, ref)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		config.Image = trustedRef.String()
	// 	}
	// }

	//create the container
	response, err := docker.Client.ContainerCreate(ctx, config, hostConfig, networkingConfig, name)
	er.CheckError(err)
	//if image not found try to pull it
	if err != nil {
		if apiclient.IsErrImageNotFound(err) {
			// fmt.Fprintf(stderr, "Unable to find image '%s' locally\n", ref.String())

			// we don't want to write to stdout anything apart from container.ID
			if err = pullImage(ctx, docker, config.Image, stderr); err != nil {
				return cont.ContainerCreateCreatedBody{}, err
			}
			// if ref, ok := ref.(reference.NamedTagged); ok != nil {
			// 	if err := docker.TagTrusted(ctx, trustedRef, ref); err != nil {
			// 	return nil, err
			// 	}
			// }
			// Retry
			var retryErr error
			response, retryErr = docker.Client.ContainerCreate(ctx, config, hostConfig, networkingConfig, name)
			if retryErr != nil {
				return cont.ContainerCreateCreatedBody{}, retryErr
			}
		} else {
			return cont.ContainerCreateCreatedBody{}, err
		}
	}

	for _, warning := range response.Warnings {
		fmt.Fprintf(stderr, "WARNING: %s\n", warning)
	}
	if containerIDFile != nil {
		if err = containerIDFile.Write(response.ID); err != nil {
			return cont.ContainerCreateCreatedBody{}, err
		}
	}
	return response, nil
}
