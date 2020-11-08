package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	ocidb "github.com/ocidb/ocidb/pkg/ocidb"
	ocidbtypes "github.com/ocidb/ocidb/pkg/ocidb/types"
)

func main() {
	port := 0
	if os.Getenv("OCIDB_PORT") != "" {
		p, err := strconv.Atoi(os.Getenv("OCIDB_PORT"))
		if err != nil {
			panic(err)
		}

		port = p
	}

	connectOpts := ocidbtypes.ConnectOpts{
		Host:      os.Getenv("OCIDB_HOST"),
		Port:      port,
		Namespace: os.Getenv("OCIDB_NAMESPACE"),
		Username:  os.Getenv("OCIDB_USERNAME"),
		Password:  os.Getenv("OCIDB_PASSWORD"),
		Database:  os.Getenv("OCIDB_DATABASE"),
	}

	fmt.Printf("%#v\n", connectOpts)
	connection, err := ocidb.Connect(context.TODO(), &connectOpts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\b", connection)
}
