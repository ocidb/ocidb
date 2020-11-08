package ocidb

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	"github.com/ocidb/ocidb/pkg/ocidb/types"
	"github.com/pkg/errors"
)

var ErrNotInitialized = errors.New("not_initialized")

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
	allowedMediaTypes := []string{"ocidb.index"}

	err := pull(ctx, resolver, indexImageRef, fileStore, allowedMediaTypes)
	if isNotInitializedErr(err) {
		fmt.Printf("need to initialized")
		return nil, nil
	}

	return &connection, nil
}

func isNotInitializedErr(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == ErrNotInitialized.Error()
}

func pull(ctx context.Context, resolver remotes.Resolver, ref string, ingester *content.FileStore, allowedMediaTypes []string) error {
	desc, _, err := oras.Pull(ctx, resolver, ref, ingester, oras.WithAllowedMediaTypes(allowedMediaTypes))
	if err != nil {
		if strings.HasSuffix(err.Error(), " not found") {
			return ErrNotInitialized
		}
		return errors.Wrap(err, "failed to pull index file")
	}

	fmt.Printf("%#v\n", desc)
	return nil
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
