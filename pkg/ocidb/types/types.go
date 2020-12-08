package types

import (
	"database/sql"

	schemasv1alpha4 "github.com/schemahero/schemahero/pkg/apis/schemas/v1alpha4"
)

type ConnectOpts struct {
	Host      string
	Port      int
	PlainHTTP bool
	Namespace string
	Username  string
	Password  string

	Database string

	Tables []schemasv1alpha4.TableSpec // A set of SchemaHero schema definitions to apply

	ReadOnly bool // when set, nothing will be commited to the database, all writes are disabled
}

type Connection struct {
	ConnectOpts   *ConnectOpts
	LocalCacheDir string
	DB            *sql.DB
}
