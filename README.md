# go-gin-wrapper

[![Build Status](https://travis-ci.org/hiromaily/go-gin-wrapper.svg?branch=master)](https://travis-ci.org/hiromaily/go-gin-wrapper)
[![Coverage Status](https://coveralls.io/repos/github/hiromaily/go-gin-wrapper/badge.svg)](https://coveralls.io/github/hiromaily/go-gin-wrapper)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-gin-wrapper)](https://goreportcard.com/report/github.com/hiromaily/go-gin-wrapper)
[![codebeat badge](https://codebeat.co/badges/30d3509a-36be-408b-bfed-ddd6f601c075)](https://codebeat.co/projects/github-com-hiromaily-go-gin-wrapper-master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/83df70f4373d47f8b889be10ad28d4a3)](https://www.codacy.com/app/hiromaily2/go-gin-wrapper?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=hiromaily/go-gin-wrapper&amp;utm_campaign=Badge_Grade)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-gin-wrapper/master/LICENSE)
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/hiromaily/go-gin-wrapper)

Go-gin-wrapper is wrapper of go gin web framework plus reverseproxy. React by ES6 is used on part of front-end. 
[gin-gonic/gin](https://github.com/gin-gonic/gin)  

This project has started since 2016 to study Golang and code is quite messy.
Now it's under `refactoring`.

## Refactoring
- [ ] change architecture like Clean Architecture
- [ ] remove any dependencies from [hiromaily/golibs](https://github.com/hiromaily/golibs)
  - [x] use [sqlboiler](https://github.com/volatiletech/sqlboiler) as ORM
  - [x] remove MongoDB
  - [x] replace log to zap logger
- [x] add zap logger 
- [ ] clean up variable name
- [ ] clean up comments
- [ ] catch up with latest [gin](https://github.com/gin-gonic/gin)
- [ ] update front-end
- [ ] unittest by table driven test
- [ ] refactoring code by [The Second Edition of "Refactoring"](https://martinfowler.com/articles/refactoring-2nd-ed.html)
- [ ] switch any Japanese to English

## Example
This is built on Heroku. You can see [here](https://ginserver.herokuapp.com/).


## Installation
```
$ go get github.com/hiromaily/go-gin-wrapper ./...
 or
$ go get github.com/hiromaily/go-gin-wrapper ./cmd/ginserver/...
```

#### Setup for local environment with Docker
```
# it would be better to execute build, up separately for first the build
$ docker-compose build
$ docker-compose up
```

## Configuration

### 1. Common settings
#### MySQL
It requires set database and table.  
If you want know more details, use docker environment usign docker-create.sh.  
This is the easiest way to configure whole environment.


#### Redis
If you store session data on Redis, Redis is required. 
But it's not indispensable.

#### TOML file
This is for configuration.
```
$ cp ./configs/settings.default.toml ./configs/settings.toml

```
When running web server, go-gin-wrapper requires toml file as configuration information.  
Without command line arguments for toml file path, this try to read ```configs/settings.toml```.   
If you want to use original toml file, use command line arguments ```-f filepath```.  
```
ginserver -f /app/config/yourfile.toml
```

* server
* proxy
* auth
* mysql
* redis
* aws  

※ As needed, secret information can be ciphered.(using AES encryption)

#### Authentication for Login
It's available using Google or Facebook account using OAuth2 authentication.

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
