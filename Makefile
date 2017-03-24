# Note: tabs by space can't not used for Makefile!
MONGO_PORT=27017

###############################################################################
# Initial Settings
###############################################################################
dcinit:
	docker-compose up mysql redis mongo

mongoinit:
	mongo 127.0.0.1:$(MONGO_PORT)/admin --eval "var port = $(MONGO_PORT);" ./docker_build/mongo/init.js
	mongorestore -h 127.0.0.1:${MONGO_PORT} --db hiromaily docker_build/mongo/dump/hiromaily


dbinit:
	export DB_NAME=hiromaily
	sh ./data/sql/setup.sh


###############################################################################
# Docker
###############################################################################
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


###############################################################################
# Managing Dependencies
###############################################################################
changegit:
	cd ${GOPATH}/src/github.com/aws/aws-sdk-go
	git checkout v0.9.17

update:
	go get -u -v ./...
	go get -u github.com/tools/godep

godep:
	rm -rf Godeps
	rm -rf ./vendor
	godep save ./...


###############################################################################
# Golang detection and formatter
###############################################################################
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


###############################################################################
# Build
###############################################################################
bld:
	go build -i -v -o ${GOPATH}/bin/ginserver ./cmd/ginserver/

bldproxy:
	go build -i -v -o ${GOPATH}/bin/reverseproxy ./cmd/reverseproxy/

bldswg:
	go build -i -v -o ${GOPATH}/bin/swgserver ./swagger/go-swagger/cmd/swagger-server/


###############################################################################
# Execution
###############################################################################
run:
	go run ./cmd/ginserver/main.go

exec:
	ginserver -f ./configs/settings.toml

execproxy:
	PORTS=(9997 9998 9999)
	for port in ${PORTS[@]}
	do
		echo 'port is ${port}'
		ginserver -f ./configs/settings.toml -P ${port} &
	done
	sleep 5s

	reverseproxy -f ./configs/settings.toml
	#proxy.hiromaily.com:9990

execswg:
	swgserver


###############################################################################
# Test
###############################################################################
quicktest:
	go test -run TestLogin -v cmd/ginserver/*.go -f ../../configs/settings.toml

test:
	# Create Test Data
	export DB_NAME=hiromaily2
	export DB_PORT=13306
	export DB_USER=root
	export DB_PASS=root
	sh ./data/sql/setup.sh

	# Execute
	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-f ../../configs/settings.toml -om 0

	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-run "TestGetUserAPIRequestOnTable" \
	-f ../../configs/settings.toml -om 1

	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-run "TestGetUserAPIRequestOnTable" \
	-f ../../configs/settings.toml -om 2

	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-run "TestGetJwtAPIRequestOnTable|TestGetUserAPIRequestOnTable" \
	-f ../../configs/settings.toml -om 1

	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-run "TestGetJwtAPIRequestOnTable|TestGetUserAPIRequestOnTable" \
	-f ../../configs/settings.toml -om 2


###############################################################################
# Heroku
###############################################################################
herokudeploy:
	git push -f heroku master
