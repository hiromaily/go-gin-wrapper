# go-gin-wrapper

[![Build Status](https://travis-ci.org/hiromaily/go-gin-wrapper.svg?branch=master)](https://travis-ci.org/hiromaily/go-gin-wrapper)
[![Coverage Status](https://coveralls.io/repos/github/hiromaily/go-gin-wrapper/badge.svg?branch=master)](https://coveralls.io/github/hiromaily/go-gin-wrapper?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-gin-wrapper)](https://goreportcard.com/report/github.com/hiromaily/go-gin-wrapper)
[![codebeat badge](https://codebeat.co/badges/30d3509a-36be-408b-bfed-ddd6f601c075)](https://codebeat.co/projects/github-com-hiromaily-go-gin-wrapper-master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/83df70f4373d47f8b889be10ad28d4a3)](https://www.codacy.com/app/hiromaily2/go-gin-wrapper?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=hiromaily/go-gin-wrapper&amp;utm_campaign=Badge_Grade)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-gin-wrapper/master/LICENSE)
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/hiromaily/go-gin-wrapper)

Go-gin-wrapper is wrapper of [gin-gonic/gin](https://github.com/gin-gonic/gin) web framework plus reverse proxy.  
React by ES6 is used on part of front-end. But it's quite outdated. 

This project has started since 2016 to study Golang and code is quite messy.
Now it's under `refactoring`.


## Refactoring
- [ ] refactoring code by [The Second Edition of "Refactoring"](https://martinfowler.com/articles/refactoring-2nd-ed.html)
- [ ] refactoring code by [codebeat.co](https://codebeat.co/projects/github-com-hiromaily-go-gin-wrapper-master/ratings)
- [ ] change architecture like Clean Architecture
  - [ ] make UI adaptable like WEB to CUI
- [ ] fix session functionality
- [x] remove any dependencies from [hiromaily/golibs](https://github.com/hiromaily/golibs)
  - [x] use [sqlboiler](https://github.com/volatiletech/sqlboiler) as ORM
  - [x] remove MongoDB
  - [x] replace log to zap logger
- [x] add zap logger
- [x] add graceful shutdown
- [x] refactoring jwt package  
- [x] fix main_test
- [x] unittest by table driven test
- [x] clean up variable, func name
- [x] clean up comments
- [ ] catch up with latest [gin](https://github.com/gin-gonic/gin)
- [ ] update front-end

## Requirements
- Golang 1.15+
- Docker compose
  - MySQL 5.7
  - Redis
- [direnv](https://github.com/direnv/direnv) for MacOS user  

## Example
Example is [here](https://ginserver.herokuapp.com/) on Heroku.


## Functionalities
### Authentication for Login
`OAuth2` authentication with Google/Facebook is available.

#### Authentication for API
`JWT(Json Web Token)` is used  for authentication.
See configuration `[api.auth]` in toml file.

#### APIs documentation by [Swagger](http://localhost:8080/swagger/?url=https://raw.githubusercontent.com/hiromaily/go-gin-wrapper/master/swagger/swagger.yaml)


## Configuration
See `./configs/settings.default.toml`  
As needed, secret information can be encrypted.(using AES encryption)

## Dependent middleware
- MySQL
- Redis whose used as session store


## Installation on local (MacOS is expected)
```
# 1. clone repository
$ git clone https://github.com/hiromaily/go-gin-wrapper.git

# 2. copy settings.default.toml
$ cp configs/settings.default.toml configs/settings.toml

# 3. edit `configs/settings.toml`
# 3.1. you may need to modify in settings.toml
 [server.docs]
 # set `go-gin-wrapper` path
 # this path must be chnaged first for specific environment
 path = "${GOPATH}/src/github.com/hiromaily/go-gin-wrapper"

$ 4. start MySQL
$ docker-compose up mysql
 
# 5. make sure `go run` works
$ go run ./cmd/ginserver/ -f ./configs/settings.toml

# 6. unit test requires ${GOGIN_CONF} environment variable
- define proper path on `.envrc`
```


## Installation on Docker
```
$ docker-compose build
$ docker-compose up
```

#### Docker related files
* configs/docker.toml
* docker-compose.yml
* docker-compose.override.yml
* Dockerfile
* ./build/docker/*


## Installation on Heroku
```
## Install 
$ heroku create ginserver --buildpack heroku/go

# MySQL
$ heroku addons:create cleardb
$ heroku config | grep CLEARDB_DATABASE_URL

# Redis
$ heroku addons:create heroku-redis:hobby-dev -a ginserver 
$ heroku config | grep REDIS

## Environment variable
$ heroku config:add ENC_KEY=xxxxx
$ heroku config:add ENC_IV=xxxxx

## Check
$ heroku config | grep CLEARDB_DATABASE_URL
$ heroku config | grep REDIS
$ heroku ps -a ginserver

## Deploy
$ git push -f heroku master
``` 

* [site on heroku](https://ginserver.herokuapp.com/)

* `./configs/heroku.toml` is used to run server.  
```
$ ginserver -f /app/configs/heroku.toml -crypto
```


## Environment variables
| NAME              | Value            | Explanation                        |
|:------------------|:-----------------|:-----------------------------------|
| GOGIN_CONF        | xxxxx            | config path, required in unit test |
| ENC_KEY           | xxxxx            | encryption                         |
| ENC_IV            | xxxxx            | encryption                         |

#### Only Heroku environment
| NAME              | Value            |
|:------------------|:-----------------|
| PORT              | 8080             |

- Heroku server use `PORT` automatically as environment variable.


## Usage
```
Usage: ginserver [options...]

Options:
  -f      Toml file path
  -p      Overwriten server port number
  -crypto if true, values in config file are encrypted

e.g.
 $ ginserver -f /app/configs/yourfile.toml -crypto
```


## Profiling
See config file.
```
[profile]
enable = true
```

After running ginserver, acccess the below links.
```
[GIN-debug] GET    /debug/pprof/
[GIN-debug] GET    /debug/pprof/heap
[GIN-debug] GET    /debug/pprof/goroutine
[GIN-debug] GET    /debug/pprof/block
[GIN-debug] GET    /debug/pprof/threadcreate
[GIN-debug] GET    /debug/pprof/cmdline
[GIN-debug] GET    /debug/pprof/profile
[GIN-debug] GET    /debug/pprof/symbol
[GIN-debug] POST   /debug/pprof/symbol
[GIN-debug] GET    /debug/pprof/trace
```
