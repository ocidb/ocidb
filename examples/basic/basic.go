package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	ocidb "github.com/ocidb/ocidb/pkg/ocidb"
	ocidbtypes "github.com/ocidb/ocidb/pkg/ocidb/types"
	schemasv1alpha4 "github.com/schemahero/schemahero/pkg/apis/schemas/v1alpha4"
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

	tables := []schemasv1alpha4.TableSpec{
		{
			Name: "testing",
			Schema: &schemasv1alpha4.TableSchema{
				SQLite: &schemasv1alpha4.SqliteTableSchema{
					Columns: []*schemasv1alpha4.SqliteTableColumn{
						{
							Name: "test",
							Type: "text",
						},
					},
				},
			},
		},
	}

	connectOpts := ocidbtypes.ConnectOpts{
		Host:      os.Getenv("OCIDB_HOST"),
		Port:      port,
		Namespace: os.Getenv("OCIDB_NAMESPACE"),
		Username:  os.Getenv("OCIDB_USERNAME"),
		Password:  os.Getenv("OCIDB_PASSWORD"),
		Database:  os.Getenv("OCIDB_DATABASE"),
		Tables:    tables,
	}

	fmt.Printf("%#v\n", connectOpts)
	connection, err := ocidb.Connect(context.TODO(), &connectOpts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\b", connection)
}
