# go-gin-wrapper

[![Build Status](https://travis-ci.org/hiromaily/go-gin-wrapper.svg?branch=master)](https://travis-ci.org/hiromaily/go-gin-wrapper)
[![Coverage Status](https://coveralls.io/repos/github/hiromaily/go-gin-wrapper/badge.svg)](https://coveralls.io/github/hiromaily/go-gin-wrapper)
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
- [ ] change architecture like Clean Architecture
- [x] remove any dependencies from [hiromaily/golibs](https://github.com/hiromaily/golibs)
  - [x] use [sqlboiler](https://github.com/volatiletech/sqlboiler) as ORM
  - [x] remove MongoDB
  - [x] replace log to zap logger
- [x] add zap logger
- [x] add graceful shutdown
- [x] refactoring jwt package  
- [ ] clean up variable, func name
- [ ] clean up comments
- [ ] catch up with latest [gin](https://github.com/gin-gonic/gin)
- [ ] update front-end
- [ ] unittest by table driven test
- [ ] refactoring code by [The Second Edition of "Refactoring"](https://martinfowler.com/articles/refactoring-2nd-ed.html)
- [ ] switch any Japanese to English
- [ ] refactoring and fix test

## Example
Example is [here](https://ginserver.herokuapp.com/) on Heroku.


## Installation
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
```


## Configuration
See `./configs/settings.default.toml`  
As needed, secret information can be encrypted.(using AES encryption)

## Dependent middleware
- MySQL 
- Redis whose used as session store



#### Authentication for Login
`OAuth2` authentication with Google/Facebook is available.

#### Authentication for API
It's implemented by JWT(Json Web Token) for authentication.
Set ```[api.auth]``` on toml file.
You can choose HMAC or RSA as signature pattern.


### 2. On heroku
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
$ heroku config:add HEROKU_FLG=1
$ heroku config:add ENC_KEY=xxxxx
$ heroku config:add ENC_IV=xxxxx

## Check
$ heroku config | grep CLEARDB_DATABASE_URL
$ heroku config | grep REDIS
$ heroku ps -a ginserver

## Deploy
$ git push -f heroku master
``` 

* Access (For check hot to work)  
[site on heroku](https://ginserver.herokuapp.com/)


* Heroku environment set configs/heroku.toml when starting to run.  
```
ginserver -f /app/configs/heroku.toml -crypto
```

### 3. On Docker
Docker environment set data/toml/docker.toml when starting to run.  

#### Docker related files
* docker-compose.yml
* docker-compose.override.yml
* Dockerfile
* ./build/docker/*


## Environment variable e.g.
### 1. Common
| NAME              | Value            |
|:------------------|:-----------------|
| ENC_KEY           | xxxxx            |
| ENC_IV            | xxxxx            |

### 2. For Heroku environment
| NAME              | Value            |
|:------------------|:-----------------|
| HEROKU_FLG        | 1                |
| PORT              | 8080             |

Heroku server use ```PORT``` automatically as environment variable.


## APIs
Documentation is prepared using Swagger
[swagger](http://localhost:8080/swagger/?url=https://raw.githubusercontent.com/hiromaily/go-gin-wrapper/master/swagger/swagger.yaml)


## Usage
```
Usage: ginserver [options...]

Options:
  -f     Toml file path

e.g.
 $ ginserver -f /app/configs/yourfile.toml
```


## Profiling
Set config first.
```
[profile]
enable = true
```

After running werserver, acccess tha below links.
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
