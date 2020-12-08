package ocidb

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ocidb/ocidb/pkg/ocidb/types"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	schemasv1alpha4 "github.com/schemahero/schemahero/pkg/apis/schemas/v1alpha4"
	"github.com/schemahero/schemahero/pkg/database"
)

var ErrNotInitialized = errors.New("not_initialized")

// Commit will push to the registry
func Commit(ctx context.Context, connection *types.Connection) error {
	// TODO set a lock on the database
	connection.DB.Close()

	data, err := ioutil.ReadFile(filepath.Join(connection.LocalCacheDir, "database.db"))
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	resolver := docker.NewResolver(docker.ResolverOptions{})
	ref := fmt.Sprintf("%s:index", imageRefFromConnectOpts(connection.ConnectOpts))

	memoryStore := content.NewMemoryStore()
	pushContents := []ocispec.Descriptor{
		memoryStore.Add("database.db", "ocidb.db", data),
	}
	desc, err := oras.Push(ctx, resolver, ref, memoryStore, pushContents)
	if err != nil {
		return errors.Wrap(err, "failed to push created sqlite database")
	}

	fmt.Printf("%#v\n", desc)

	db, err := sql.Open("sqlite3", filepath.Join(connection.LocalCacheDir, "database.db"))
	if err != nil {
		return errors.Wrap(err, "failed to open")
	}
	connection.DB = db

	return nil
}

// Connect is called by the application to create a connection to an existing
// database. The registry details are reuqired, along with the database name.
// All other parameters are optional as they have sane defaults.
func Connect(ctx context.Context, connectOpts *types.ConnectOpts) (*types.Connection, error) {
	connection := types.Connection{
		ConnectOpts:   connectOpts,
		LocalCacheDir: os.TempDir(),
	}

	resolver := docker.NewResolver(docker.ResolverOptions{
		PlainHTTP: connectOpts.PlainHTTP,
	})

	indexImageRef := fmt.Sprintf("%s:index", imageRefFromConnectOpts(connectOpts))

	fileStore := content.NewFileStore(connection.LocalCacheDir)
	defer fileStore.Close()
	allowedMediaTypes := []string{"ocidb.db"}

	desc := &ocispec.Descriptor{}
	var err error
	desc, err = pull(ctx, resolver, indexImageRef, fileStore, allowedMediaTypes)
	if isNotInitializedErr(err) {
		if err := initialize(ctx, connectOpts.Database, resolver, indexImageRef, connectOpts.Tables); err != nil {
			return nil, errors.Wrap(err, "failed to initialize new db")
		}

		desc, err = pull(ctx, resolver, indexImageRef, fileStore, allowedMediaTypes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to pull newly initialized db")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to pull database")
	}

	fmt.Printf("%s\n", desc.Digest)

	db, err := sql.Open("sqlite3", filepath.Join(connection.LocalCacheDir, "database.db"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sqlite database")
	}

	connection.DB = db

	return &connection, nil
}

func isNotInitializedErr(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == ErrNotInitialized.Error()
}

func initialize(ctx context.Context, databaseName string, resolver remotes.Resolver, ref string, tables []schemasv1alpha4.TableSpec) error {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return errors.Wrap(err, "failed to create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	localPath := filepath.Join(tmpDir, fmt.Sprintf("%s.db", databaseName))
	db, err := sql.Open("sqlite3", localPath)
	if err != nil {
		return errors.Wrap(err, "failed to open")
	}
	db.Close()

	// schemahero applies the schemas
	schemaheroDatabase := database.Database{
		Driver: "sqlite",
		URI:    localPath,
	}
	for _, table := range tables {
		statements, err := schemaheroDatabase.PlanSyncTableSpec(&table)
		if err != nil {
			return errors.Wrap(err, "failed to plan schema migration")
		}

		if err := schemaheroDatabase.ApplySync(statements); err != nil {
			return errors.Wrap(err, "failed to apply statements")
		}
	}

	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		return errors.Wrap(err, "failed to read created sqlite file")
	}

	memoryStore := content.NewMemoryStore()
	pushContents := []ocispec.Descriptor{
		memoryStore.Add("database.db", "ocidb.db", data),
	}
	desc, err := oras.Push(ctx, resolver, ref, memoryStore, pushContents)
	if err != nil {
		return errors.Wrap(err, "failed to push created sqlite database")
	}

	fmt.Printf("%#v\n", desc)
	return nil
}

func pull(ctx context.Context, resolver remotes.Resolver, ref string, ingester *content.FileStore, allowedMediaTypes []string) (*ocispec.Descriptor, error) {
	desc, _, err := oras.Pull(ctx, resolver, ref, ingester, oras.WithAllowedMediaTypes(allowedMediaTypes))
	if err != nil {
		if strings.HasSuffix(err.Error(), " not found") {
			return nil, ErrNotInitialized
		}
		return nil, errors.Wrap(err, "failed to pull index file")
	}

	return &desc, nil
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
