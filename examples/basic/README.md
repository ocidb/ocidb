# Basic example

This example shows a simple implementation of a go function that writes and then queries an oci db.

To run:

```
export OCIDB_HOST=index.docker.io
export OCIDB_USERNAME=username
export OCIDB_PASSWORD=token
export OCIDB_NAMESPACE=username
export OCIDB_DATABASE=ocidb-basic

go run ./basic.go
```