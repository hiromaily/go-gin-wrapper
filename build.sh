#!/bin/sh

###########################################################
# Variable
###########################################################
#export GOTRACEBACK=single
GOTRACEBACK=all
CURRENTDIR=`pwd`

TEST_MODE=1
AUTO_EXEC=0
GODEP_MODE=1
AUTO_GITCOMMIT=1
HEROKU_MODE=1
DOCKER_MODE=0

###########################################################
# Update all package
###########################################################
#go get -u -v ./...
#go get -d -v ./...
#go get -u github.com/tools/godep


###########################################################
# Adjust version dependency of projects
###########################################################
#cd ${GOPATH}/src/github.com/aws/aws-sdk-go
#git checkout v0.9.17
#git checkout master


###########################################################
# go fmt and go vet
###########################################################
echo '============== go fmt; go vet; =============='
go fmt ./...
#go vet ./...
go vet `go list ./... | grep -v '/vendor/'`
EXIT_STATUS=$?

if [ $EXIT_STATUS -gt 0 ]; then
    exit $EXIT_STATUS
fi

###########################################################
# go lint
###########################################################
# it's too strict
#golint ./...


###########################################################
# go list for check import package
###########################################################
#go list -f '{{.ImportPath}} -> {{join .Imports "\n"}}' ./cmd/ginserver/main.go


###########################################################
# go build and install
###########################################################
echo '============== go build -i -v -o; =============='
rm -rf Godeps
rm -rf ./vendor

#-n show just command for build
#go build -i -n ./cmd/ginserver/

#rebuild dependent packages (rebuild all package)
#go build -a -v -o ${GOPATH}/bin/ginserver ./cmd/ginserver/

#build and install
go build -i -v -o ${GOPATH}/bin/ginserver ./cmd/ginserver/
EXIT_STATUS=$?

if [ $EXIT_STATUS -gt 0 ]; then
    exit $EXIT_STATUS
fi


###########################################################
# test
###########################################################
if [ $TEST_MODE -eq 1 ]; then
    echo '============== test =============='

    go test -v cmd/ginserver/*.go -f ../../configs/settings.toml
    EXIT_STATUS=$?
    echo EXIT_STATUS

fi
#stress test
#https://github.com/rakyll/boom
# $ boom -n 1000 -c 100 https://google.com


###########################################################
# exec
###########################################################
if [ $AUTO_EXEC -eq 1 ]; then
    echo '============== exec =============='
    if [ $HEROKU_MODE -eq 1 ]; then
        #HEROKU ENV
        export HEROKU_FLG=1
        #export CLEARDB_DATABASE_URL=mysql://be2ebea7cda583:49eef93c@us-cdbr-iron-east-04.cleardb.net/heroku_aa95a7f43af0311?reconnect=true
        #export REDIS_URL=redis://h:pf217irr4gts39d29o0012bghsi@ec2-50-19-83-130.compute-1.amazonaws.com:20819

        ginserver -f ./configs/heroku.toml
    else
        #ginserver -f ${PWD}/configs/settings.toml
        ginserver -f ./configs/settings.toml
    fi
fi

###########################################################
# cross-compile for linux
###########################################################
#GOOS=linux go install -v ./...


###########################################################
# godep
###########################################################
if [ $GODEP_MODE -eq 1 ]; then
    echo '============== godeps =============='

    #go get -u github.com/tools/godep

    #Save
    rm -rf Godeps
    rm -rf vendor

    godep save ./...
    EXIT_STATUS=$?

    if [ $EXIT_STATUS -gt 0 ]; then
        exit $EXIT_STATUS
    fi
fi

#Build
#godep go build -o book ./cmd/book/

#Restore
#godep restore


###########################################################
# git add, commit, push
###########################################################
if [ $AUTO_GITCOMMIT -eq 1 ]; then
    echo '============== git recm, pufom =============='
    git recm
    git pufom
    git st
fi


###########################################################
# heroku
###########################################################
if [ $HEROKU_MODE -eq 1 ]; then
    echo '============== heroku: git push =============='
    git push -f heroku master
fi

#heroku config:add HEROKU_FLG=1

#heroku ps -a ginserver
#heroku run bash
#heroku logs -t
#heroku ps
#heroku config

#heroku open

###########################################################
# endpoint
###########################################################
#Local
#http://hiromaily.com:9999

#Heroku
#https://ginserver.herokuapp.com/


###########################################################
# Docker
###########################################################
if [ $DOCKER_MODE -eq 1 ]; then
    echo '============== docker =============='
    ./docker-create.sh
fi


###########################################################
# godoc
###########################################################
#godoc -http :8000
#http://localhost:8000/pkg/

