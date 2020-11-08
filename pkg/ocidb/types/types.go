package types

type ConnectOpts struct {
	Host      string
	Port      int
	Namespace string
	Username  string
	Password  string

	Database string

	ReadOnly bool // when set, nothing will be commited to the database, all writes are disabled
}

type Connection struct {
	LocalCacheDir string
}
