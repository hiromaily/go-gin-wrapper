# Note: tabs by space can't not used for Makefile!
MONGO_PORT=27017

dcinit:
	docker-compose up mysql redis mongo

setupmongo:
	mongo 127.0.0.1:$(MONGO_PORT)/admin --eval "var port = $(MONGO_PORT);" ./docker_build/mongo/init.js
	mongorestore -h 127.0.0.1:${MONGO_PORT} --db hiromaily docker_build/mongo/dump/hiromaily

dcstart:
	docker start web-redisd
	docker start web-mysqld
	docker start web-mongod

dcstop:
	docker stop web-redisd
	docker stop web-mysqld
	docker stop web-mongod

dcmongo:
	docker exec -it web-mongo bash

update:
	rm -rf Godeps
	rm -rf ./vendor
	go get -u -v ./...
	go get -u github.com/tools/godep
	godep save ./...

fmt:
	go fmt `go list ./... | grep -v '/vendor/'`

vet:
	go vet `go list ./... | grep -v '/vendor/'`

lint:
	golint ./... | grep -v '^vendor\/' || true
	misspell `find . -name "*.go" | grep -v '/vendor/'`
	ineffassign .

chk:
	go fmt `go list ./... | grep -v '/vendor/'`
	go vet `go list ./... | grep -v '/vendor/'`
	golint ./... | grep -v '^vendor\/' || true
	misspell `find . -name "*.go" | grep -v '/vendor/'`
	ineffassign .

run:
	go run ./cmd/ginserver/main.go

bld:
	go build -i -v -o ${GOPATH}/bin/ginserver ./cmd/ginserver/

bldswg:
	go build -i -v -o ${GOPATH}/bin/swgserver ./swagger/go-swagger/cmd/swagger-server/

exec:
	ginserver -f ./configs/settings.toml

execswg:
	swgserver

godep:
	echo go-gin was modified by me because of bug. So it could not be worked.
	#go get -u github.com/tools/godep
	#rm -rf Godeps
	#rm -rf ./vendor
	#godep save ./...

#develop: pull
#	docker-compose build
#	cd ../sendy; make develop; make build
#	cd ../backoffice; make develop
#	cd ../golang/src/hugoevents/goapi; GOPATH=$(PWD)/../golang glide install
