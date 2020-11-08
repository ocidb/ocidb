package main

import (
	"context"
	"fmt"
	"os"

	ocidb "github.com/ocidb/ocidb/pkg/ocidb"
	ocidbtypes "github.com/ocidb/ocidb/pkg/ocidb/types"
)

func main() {
	connectOpts := ocidbtypes.ConnectOpts{
		Host:      os.Getenv("OCIDB_HOST"),
		Namespace: os.Getenv("OCIDB_NAMESPACE"),
		Username:  os.Getenv("OCIDB_USERNAME"),
		Password:  os.Getenv("OCIDB_PASSWORD"),
		Database:  os.Getenv("OCIDB_DATABASE"),
	}

	connection, err := ocidb.Connect(context.TODO(), &connectOpts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\b", connection)
}
