#!/bin/sh
###
# initialize for docker environment
###

go get -d -v ./...
go build -v -o /go/bin/ginserver ./cmd/ginserver/

ginserver -f ./configs/docker.toml
#it was impossible to passs environment variable from build.sh
#if [ $RUN_TEST -eq 1 ]; then
#    # Test
#    go test -v cmd/ginserver/*.go -f ../../configs/docker.toml
#else
#    # Run web server
#    ginserver -f ./configs/docker.toml
#fi
