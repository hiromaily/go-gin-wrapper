#!/bin/sh
###
# initialize for docker environment
###

go get -d -v ./...
go build -v -o /go/bin/ginserver ./cmd/ginserver/

ginserver -f ./configs/docker.toml