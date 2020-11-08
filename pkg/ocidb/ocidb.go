package ocidb

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd/remotes/docker"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	"github.com/ocidb/ocidb/pkg/ocidb/types"
	"github.com/pkg/errors"
)

// Connect is called by the application to create a connection to an existing
// database. The registry details are reuqired, along with the database name.
// All other parameters are optional as they have sane defaults.
func Connect(ctx context.Context, connectOpts *types.ConnectOpts) (*types.Connection, error) {
	connection := types.Connection{
		LocalCacheDir: os.TempDir(),
	}

	resolver := docker.NewResolver(docker.ResolverOptions{})

	indexImageRef := fmt.Sprintf("%s:index", imageRefFromConnectOpts(connectOpts))

	fileStore := content.NewFileStore(connection.LocalCacheDir)
	defer fileStore.Close()
	allowedMediaTypes := []string{"aaa"}

	desc, _, err := oras.Pull(ctx, resolver, indexImageRef, fileStore, oras.WithAllowedMediaTypes(allowedMediaTypes))
	if err != nil {
		return nil, errors.Wrap(err, "failed to pull index file")
	}

	fmt.Printf("desc: %#v\n", desc)

	return &connection, nil
}

func imageRefFromConnectOpts(connectOpts *types.ConnectOpts) string {
	namespace := connectOpts.Namespace
	if namespace != "" {
		namespace = fmt.Sprintf("%s/", connectOpts.Namespace)
	}

	port := connectOpts.Port
	if port == 0 {
		port = 443
	}

	return fmt.Sprintf("%s:%d/%s%s", connectOpts.Host, port, namespace, connectOpts.Database)
}
