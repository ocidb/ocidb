# Basic example

This example shows a simple implementation of a go function that writes and then queries an oci db.

To run:

1. Start a local docker distribution registry

```
docker run -it -d --rm -p 5000:5000 registry
```

2. Execute the example

```
export OCIDB_HOST=localhost
export OCIDB_PORT=5000
# export OCIDB_USERNAME=username
# export OCIDB_PASSWORD=token
export OCIDB_NAMESPACE=my-app
export OCIDB_DATABASE=ocidb-basic

go run ./basic.go
```
