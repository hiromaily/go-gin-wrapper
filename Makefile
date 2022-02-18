# Note: tabs by space can't not used on Makefile!
CURRENTDIR=`pwd`
modVer=$(shell cat go.mod | head -n 3 | tail -n 1 | awk '{print $2}' | cut -d'.' -f2)
currentVer=$(shell go version | awk '{print $3}' | sed -e "s/go//" | cut -d'.' -f2)

PROJECT_ROOT=${GOPATH}/src/github.com/hiromaily/go-gin-wrapper

###############################################################################
# setup
###############################################################################
#.PHONY: gen-jwt-key
#gen-jwt-key:
#	openssl genrsa -out private.pem -aes256 4096
#	openssl rsa -pubout -in private.pem -out public.pem

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
# dependencies
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
	go get -u github.com/volatiletech/sqlboiler
	GO111MODULE=off go get -u github.com/oxequa/realize
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u -d -v ./...


###############################################################################
# linter and formatter
###############################################################################
.PHONY: lint-all
lint-all: imports lint

.PHONY: imports
imports:
	./scripts/imports.sh

.PHONY: lint
lint:
	golangci-lint run --fix

###############################################################################
# local build
###############################################################################
.PHONY: build
build:
	go build -v -o ${GOPATH}/bin/ginserver ./cmd/ginserver/

.PHONY: build-proxy
build-proxy:
	go build -v -o ${GOPATH}/bin/reverseproxy ./cmd/reverseproxy/

.PHONY: build-swg
build-swg:
	go build -v -o ${GOPATH}/bin/swgserver ./swagger/go-swagger/cmd/swagger-server/


###############################################################################
# execution
###############################################################################
.PHONY: run
run:
	#go run ./cmd/ginserver/ -f ./configs/settings.toml -crypto
	go run ./cmd/ginserver/ -f ./configs/settings.toml

.PHONY: exec
exec:
	#ginserver -f ./configs/settings.toml -crypto
	ginserver -f ./configs/settings.toml

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
.PHONY: clean-test-cache
clean-test-cache:
	go clean -testcache

.PHONY: setup-testdb
setup-testdb:
	./scripts/create-test-db-docker.sh

.PHONY: test
test:
	go test -race -v ./...

.PHONY: integration-test
integration-test: setup-testdb
	go test -race -tags=integration -v ./...


.PHONY: user-db-test
user-db-test: setup-testdb
	go test -race -tags=integration -v ./pkg/repository/...

.PHONY: get-req-test
get-req-test: setup-testdb
	go test -race -tags=integration -run TestGetRequest -v ./cmd/ginserver/...

.PHONY: login-req-test
login-req-test: setup-testdb
	go test -race -tags=integration -run TestLoginRequest -v ./cmd/ginserver/...

.PHONY: jwt-req-test
jwt-req-test: setup-testdb
	go test -race -tags=integration -run TestJWTAPIRequest -v ./cmd/ginserver/...

.PHONY: jwt-req-test
user-req-test: setup-testdb
	go test -race -tags=integration -run TestGetUserAPIRequest -v ./cmd/ginserver/...


#.PHONY: setup-testdb
#setup-testdb:
#	# create test data
#	export DB_NAME=gogin-test &&\
#	export DB_PORT=13306 &&\
#	export DB_USER=root &&\
#	export DB_PASS=root &&\
#	sh ./scripts/create-test-db.sh

###############################################################################
# httpie
#
#[GIN-debug] GET    /                         --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.BaseIndexAction-fm (7 handlers)
#[GIN-debug] GET    /index                    --> github.com/hiromaily/go-gin-wrapper/pkg/server.(*server).setBaseRouter.func1 (7 handlers)
#[GIN-debug] HEAD   /                         --> github.com/hiromaily/go-gin-wrapper/pkg/server.(*server).setBaseRouter.func2 (7 handlers)
#[GIN-debug] GET    /login                    --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.BaseLoginGetAction-fm (7 handlers)
#[GIN-debug] POST   /login                    --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.BaseLoginPostAction-fm (9 handlers)
#[GIN-debug] PUT    /logout                   --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.BaseLogoutPutAction-fm (7 handlers)
#[GIN-debug] POST   /logout                   --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.BaseLogoutPostAction-fm (7 handlers)
#[GIN-debug] GET    /apilist/                 --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIListIndexAction-fm (7 handlers)
#[GIN-debug] GET    /apilist/index2           --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIListIndex2Action-fm (7 handlers)
#[GIN-debug] GET    /oauth2/google/signin     --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.OAuth2SignInGoogleAction-fm (7 handlers)
#[GIN-debug] GET    /oauth2/google/callback   --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.OAuth2CallbackGoogleAction-fm (7 handlers)
#[GIN-debug] GET    /oauth2/facebook/signin   --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.OAuth2SignInFacebookAction-fm (7 handlers)
#[GIN-debug] GET    /oauth2/facebook/callback --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.OAuth2CallbackFacebookAction-fm (7 handlers)
#[GIN-debug] GET    /accounts/                --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.AccountIndexAction-fm (7 handlers)
#[GIN-debug] GET    /admin/                   --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.AdminIndexAction-fm (8 handlers)
#[GIN-debug] GET    /admin/index              --> github.com/hiromaily/go-gin-wrapper/pkg/server.(*server).setAdminRouter.func1 (8 handlers)
#[GIN-debug] POST   /api/jwts                 --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIJWTIndexPostAction-fm (9 handlers)
#[GIN-debug] GET    /api/users                --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIUserListGetAction-fm (12 handlers)
#[GIN-debug] POST   /api/users                --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIUserInsertPostAction-fm (12 handlers)
#[GIN-debug] GET    /api/users/id/:id         --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIUserGetAction-fm (12 handlers)
#[GIN-debug] PUT    /api/users/id/:id         --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIUserPutAction-fm (12 handlers)
#[GIN-debug] DELETE /api/users/id/:id         --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIUserDeleteAction-fm (12 handlers)
#[GIN-debug] GET    /api/users/ids            --> github.com/hiromaily/go-gin-wrapper/pkg/server/controller.Controller.APIUserIDsGetAction-fm (12 handlers)
###############################################################################
.PHONY: check-all-http
check-all-http:
	# httpie #brew install httpie
	# jq     #brew install jq

	# httpie cheatsheet
	# https://devhints.io/httpie

	# top page
	http localhost:8080/
	http localhost:8080/index
	http head localhost:8080/

	# login
	http localhost:8080/login
	http --body --form POST http://localhost:8080/login inputEmail=foobar@gogin.com inputPassword=password
	# => "Referer is invalid"
	http --body --form POST http://localhost:8080/login Referer:http://localhost:8080/login inputEmail=foobar@gogin.com inputPassword=password
	# => "given token is empty"

	# get token `<input type="hidden" name="gintoken" value="969c72910db079ece758f4acfecd05e7">`
	# http localhost:8080/login | grep gintoken | awk '{print $4}' | sed 's/^.*"\(.*\)".*$/\1/'
	$(eval TOKEN := $(shell http --session=go-web-ginserver localhost:8080/login | grep gintoken | awk '{print $4}' | sed 's/^.*"\(.*\)".*$$/\1/'))
	http --body --form POST http://localhost:8080/login Referer:http://localhost:8080/login inputEmail=foobar@gogin.com inputPassword=password gintoken=$(TOKEN)

.PHONY: check-login
check-login:
	$(eval TOKEN := $(shell http --session=go-web-ginserver localhost:8080/login | grep gintoken | awk '{print $4}' | sed 's/^.*"\(.*\)".*$$/\1/'))
	http --body --form --session=go-web-ginserver POST http://localhost:8080/login Referer:http://localhost:8080/login inputEmail=foobar@gogin.com inputPassword=password gintoken=$(TOKEN)

.PHONY: check-token
check-token:
	$(eval TOKEN := $(shell http localhost:8080/login | grep gintoken | awk '{print $4}' | sed 's/^.*"\(.*\)".*$$/\1/'))
	echo $(TOKEN)


###############################################################################
# CI
###############################################################################
.PHONY: check-ci-compose
check-ci-compose:
	docker-compose -f docker-compose.ci.yml build mysql
	docker-compose -f docker-compose.ci.yml up mysql

# Note: change `GO_GIN_CONF` in `.envrc`
# then run `make integration-test`


###############################################################################
# Docker-Compose
###############################################################################

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
# Tools
# Note: environment variable `ENC_KEY`, `ENC_IV` should be set in advance
###############################################################################
.PHONY: tool-encode
tool-encode:
	go run ./tools/encryption/ -encode important-password

.PHONY: tool-decode
tool-decode:
	go run ./tools/encryption/ -decode o5PDC2aLqoYxhY9+mL0W/IdG+rTTH0FWPUT4u1XBzko=

.PHONY: tool-md5
tool-md5:
	go run ./tools/hash/ -md5 -salt1 foo-bar -salt2 hoge-hoge -target password
	go run ./tools/hash/ -md5 -salt1 foo-bar -salt2 hoge-hoge -target secret-string

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
