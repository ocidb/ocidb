
SHELL := /bin/bash

export GO111MODULE=on

all: fmt vet test

.PHONY: test
test: fmt vet
	go test ./pkg/... -coverprofile cover.out

.PHONY: fmt
fmt:
	go fmt ./pkg/... 

.PHONY: vet
vet:
	go vet ./pkg/... 
