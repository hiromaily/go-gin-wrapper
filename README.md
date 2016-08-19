# go-gin-wrapper
[![Build Status](https://travis-ci.org/hiromaily/go-gin-wrapper.svg?branch=master)](https://travis-ci.org/hiromaily/go-gin-wrapper)

Go-gin-wrapper is wrapper of go gin web framework.  
 [gin-gonic/gin](https://github.com/gin-gonic/gin)


## Installation
```
$ go get github.com/hiromaily/go-gin-wrapper ./...
```

#### For docker environment
```
$ ./docker-create.sh
```


## Configration

### 1. Common settings
#### MySQL
It requires set database and table.  
If you want know more details, use docker environment usign docker-create.sh.  
This is the easiest way to configure whole environment.


#### Redis
If you store session data on Redis, Redis is required. 
But it's not indispensable.

#### TOML file
```
$ cp configs/settings.default.toml configs/settings.toml

```
When running web server, go-gin-wrapper requires toml file as configuration information.  
Without command line arguments for toml file path, this try to read ```configs/settings.toml```.   
If you want to use original toml file, use command line arguments ```-f filepath```.  
```
ginserver -f /app/configs/yourfile.toml
```


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

# MongoDB
* attach mongodb from news-mongo app on heroku because of sharing.
* it's better to attach from dashboard.

## Environment variable
$ heroku config:add HEROKU_FLG=1

## Check
$ heroku config | grep CLEARDB_DATABASE_URL
$ heroku config | grep REDIS
$ heroku ps -a ginserver

## Deploy
$ git push -f heroku master

## Access (For check hot to work)
[site on heroku](https://ginserver.herokuapp.com/)

``` 

Heroku environment set configs/heroku.toml when starting to run.  
```
ginserver -f /app/configs/heroku.toml
```

### 3. On Docker
Docker environment set configs/docker.toml when starting to run.  

#### Docker related files
* docker-create.sh
* docker-compose.yml
* docker-entrypoint.sh
* Dockerfile
* ./docker_build/*


## Environment valuable e.g.
### 1. For Heroku environment
| NAME              | Value            |
|:------------------|:-----------------|
| HEROKU_FLG        | 1                |
| PORT              | 9999             |

Heroku server use ```PORT``` automatically as environment valuable.



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