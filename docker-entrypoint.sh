#!/bin/sh
###
# initialize for docker environment
###

go get -d -v ./...
go build -v -o /go/bin/ginserver ./cmd/ginserver/

# Run web server
ginserver -f ./configs/docker.toml

# Test
#go test -v cmd/ginserver/*.go -f ../../configs/docker.toml
