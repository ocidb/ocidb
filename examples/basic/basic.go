package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ocidb/ocidb/pkg/ocidb"
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
			Name: "users",
			Schema: &schemasv1alpha4.TableSchema{
				SQLite: &schemasv1alpha4.SqliteTableSchema{
					Columns: []*schemasv1alpha4.SqliteTableColumn{
						{
							Name: "id",
							Type: "int",
						},
						{
							Name: "month",
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

	connection, err := ocidb.Connect(context.TODO(), &connectOpts)
	if err != nil {
		panic(err)
	}

	if _, err := connection.DB.Exec("insert into users (id, month) values (?, ?)", 1, "oct"); err != nil {
		panic(err)
	}

	if err := ocidb.Commit(context.TODO(), connection); err != nil {
		panic(err)
	}

	rows, err := connection.DB.Query("select id, month from users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var month string
		if err := rows.Scan(&id, &month); err != nil {
			panic(err)
		}

		fmt.Printf("%d = %s\n", id, month)
	}
}
