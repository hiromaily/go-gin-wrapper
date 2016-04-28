#!/usr/bin/env bash
#go get -u -v ./...
go fmt ./...
go vet ./...

GOTRACEBACK=all

#cd ${GOPATH}/src/github.com/aws/aws-sdk-go
#git checkout v0.9.17

#build
go build -o ./check ./check.go
#go build -i -o ./ginserver ./ginserver.go
