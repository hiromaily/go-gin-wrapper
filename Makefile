# Note: tabs by space can't not used for Makefile!
MONGO_PORT=27017
CURRENTDIR=`pwd`
modVer=$(shell cat go.mod | head -n 3 | tail -n 1 | awk '{print $2}' | cut -d'.' -f2)
currentVer=$(shell go version | awk '{print $3}' | sed -e "s/go//" | cut -d'.' -f2)

###############################################################################
# Setup
###############################################################################
#.PHONY: install-sqlboiler
#install-sqlboiler: SQLBOILER_VERSION=4.4.0
#install-sqlboiler:
#	echo SQLBOILER_VERSION is $(SQLBOILER_VERSION)
#	go get github.com/volatiletech/sqlboiler@v$(SQLBOILER_VERSION)
#	go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql@v$(SQLBOILER_VERSION)

.PHONY: sqlboiler
sqlboiler:
	sqlboiler --wipe mysql

###############################################################################
# Managing Dependencies
###############################################################################
.PHONY: check-ver
check-ver:
	#echo $(modVer)
	#echo $(currentVer)
	@if [ ${currentVer} -lt ${modVer} ]; then\
		echo go version ${modVer}++ is required but your go version is ${currentVer};\
	fi


.PHONY: update
update:
	GO111MODULE=off go get -u github.com/oxequa/realize
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u -d -v ./...


###############################################################################
# Golang formatter and detection
###############################################################################
.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: imports
imports:
	./scripts/imports.sh


###############################################################################
# Local Build
###############################################################################
.PHONY: build
build:
	go build -i -v -o ${GOPATH}/bin/ginserver ./cmd/ginserver/

.PHONY: build-proxy
build-proxy:
	go build -i -v -o ${GOPATH}/bin/reverseproxy ./cmd/reverseproxy/

.PHONY: build-swg
build-swg:
	go build -i -v -o ${GOPATH}/bin/swgserver ./swagger/go-swagger/cmd/swagger-server/


###############################################################################
# Execution
###############################################################################
.PHONY: run
run:
	go run ./cmd/ginserver/main.go -f ./configs/settings.toml -crypto

.PHONY: exec
exec:
	ginserver -f ./configs/settings.toml -crypto

.PHONY: exec-proxy
exec-proxy:
	PORTS=(8080 8081 8082)
	for port in ${PORTS[@]}
	do
		echo 'port is ${port}'
		ginserver -f ./configs/settings.toml -P ${port} -crypto &
	done
	sleep 5s

	reverseproxy -f ./configs/settings.toml
	#proxy.hiromaily.com:9990

.PHONY: exec-swg
exec-swg:
	swgserver


.PHONY: health-check
health-check:
	while :
	do
		#000 or 200 or 404
		HTTP_STATUS=`curl -LI localhost:8080/ -w '%{http_code}\n' -s -o /dev/null`
		echo $HTTP_STATUS
		if [ $HTTP_STATUS -eq 000 ]; then
			sleep 1s
		else
			docker logs web
			break
		fi
	done


###############################################################################
# Test
###############################################################################

.PHONY: test-setup
test-setup:
	# Create Test Data
	export DB_NAME=hiromaily2 &&\
	export DB_PORT=13306 &&\
	export DB_USER=root &&\
	export DB_PASS=root &&\
	sh ./scripts/create-test-db.sh

.PHONY: test
test:
	# Execute
	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-f ../../configs/settings.toml -crypto -om 0

	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-run "TestGetUserAPIRequestOnTable" \
	-f ../../configs/settings.toml -crypto -om 1

	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-run "TestGetUserAPIRequestOnTable" \
	-f ../../configs/settings.toml -crypto -om 2

	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-run "TestGetJwtAPIRequestOnTable|TestGetUserAPIRequestOnTable" \
	-f ../../configs/settings.toml -crypto -om 1

	go test -v -covermode=count -coverprofile=profile.cov cmd/ginserver/*.go \
	-run "TestGetJwtAPIRequestOnTable|TestGetUserAPIRequestOnTable" \
	-f ../../configs/settings.toml -crypto -om 2

.PHONY: test-quick
test-quick:
	go test -run TestLogin -v cmd/ginserver/*.go -f ../../configs/settings.toml -crypto


###############################################################################
# Docker-Compose
###############################################################################

.PHONY: dc-bld
dc-bld:
	docker-compose build

.PHONY: dc-up
dc-up:
	docker-compose up

.PHONY: dc-bld-up
dc-bld-up:
	docker-compose up --build


# .PHONY: dc-test
# dc-test:
# 	export RUN_TEST=1
# 	sh ./scripts/docker-create.sh

# .PHONY: dc-shell
# dc-shell:
# 	echo '============== docker =============='
# 	# create docker container
# 	export RUN_TEST=0
# 	sh ./scripts/docker-create.sh
#
# 	#wait to be ready or not.
# 	echo 'building now. it may be take over 40s.'
# 	sleep 30s
# 	while :
# 	do
# 		#000 or 200 or 404
# 		HTTP_STATUS=`curl -LI localhost:8888/ -w '%{http_code}\n' -s -o /dev/null`
# 		echo $HTTP_STATUS
# 		if [ $HTTP_STATUS -eq 000 ]; then
# 			sleep 1s
# 		else
# 			docker logs web
# 			break
# 		fi
# 	done


###############################################################################
# Create Data
###############################################################################
.PHONY: init-db
init-db:
	export DB_NAME=hiromaily
	sh ./scripts/create-test-db.sh

# .PHONY: init-mongo
# init-mongo:
# 	#After running mongodb
# 	mongo 127.0.0.1:$(MONGO_PORT)/admin --eval "var port = $(MONGO_PORT);" ./docker/mongo/init.js
# 	mongorestore -h 127.0.0.1:${MONGO_PORT} --db hiromaily docker/mongo/dump/hiromaily


###############################################################################
# Tools
# Note: environment variable `ENC_KEY`, `ENC_IV` should be set in advance
###############################################################################
.PHONY: tool-encode
tool-encode:
	go run ./tools/encryption/ -m e important-password

.PHONY: tool-decode
tool-decode:
	go run ./tools/encryption/ -m d o5PDC2aLqoYxhY9+mL0W/IdG+rTTH0FWPUT4u1XBzko=


###############################################################################
# Front End
###############################################################################
# .PHONY: deploy-js
# deploy-js:
# 	cp /Users/hy/work/go/src/github.com/hiromaily/go-gin-wrapper/frontend_workspace/react/app/dist/apilist.bundle.js \
# 	/Users/hy/work/go/src/github.com/hiromaily/go-gin-wrapper/statics/js/


###############################################################################
# Heroku
#
#heroku ps -a ginserver
#heroku run bash
#heroku logs -t
#heroku ps
#heroku config
#
#heroku open
# https://ginserver.herokuapp.com/
#
###############################################################################
.PHONY: heroku-deploy
heroku-deploy:
	git push -f heroku master

###### e.g. command for heroku #####
#heroku config:add HEROKU_FLG=1
