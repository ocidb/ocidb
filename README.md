# OCI DB

This is an experimental project to create a relational database backed by an [OCI archive](https://opencontainers.org/).

## Use Cases

We built this to solve a specific problem. Our use case was an small (5-10 tables, each with < 20 rows), infrequently changed database that was included as part of an on-prem Kubernetes application. We had Postgres which was working fine, but very much overkill for the problem. We wanted to eliminate the Postgres requirement and create a path forward to scale ALL pods to > 1 replica. A high traffic service would move to a clou-native database solution such as CockroachDB or YugabyteDB, but that was the wrong direction for our use. We wanted to eliminate the database, not require more memory and disk and requirements. Overall size of our application mattered and we wanted to invest in solving this problem.

We are _not_ using ocidb in our product yet. This is very much experimental and built to explore this option.

Any feedback is welcome!

## Supported Registries

This project will support and OCI Artifact registry.

Currently, this includes:

- [ ] DockerHub (waiting on https://github.com/docker/roadmap/issues/135)
- [x] Amazon ECR
- [x] Azure ACR
- [x] Google Artifact Registry
- [x] Docker Distribution

Unknown support (needs research):

- [ ] GitLab Registry
- [ ] GitHub Container Registry
- [ ] Nexus Sonatype
- [ ] Artifactory

## Design Goals

**Embedded**  

OCI DB should be embedded in an existing Go application.

**Distributed**  

It should be possible to run more than 1 replica of the Go application and have some guarantees about the consistency of the data.

**Compatibility**  

This should be compatible with any OCI compatible image registry (DockerHub, ECR, GCR, etc).

**Small and Lightweight**  

We don't want extra storage requirements or memory added.

**User-mode**  

OCI DB should not require any additional (cluster-admin) access to the Kubernetes cluster.

**Reasonably fast and efficient**  

Speed of queries and efficiency of bandwidth is not the primary design goal but a secondary one. The is our open question, can we make this usable? Every query (probably) cannot go back to an OCI registry and be performant.

**Transactions**  

There should be support for transactions across multiple tables.
