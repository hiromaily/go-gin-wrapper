# Dcokerfile for go-gin-wrapper

FROM golang:1.6

#ARG redisHostName=default-redis-server
#ARG mysqlHostName=default-mysql-server

#COPY ./go-gin-wrapper /go/src/github.com/hiromaily/go-gin-wrapper/
RUN mkdir -p /go/src/github.com/hiromaily/go-gin-wrapper
RUN mkdir /var/log/goweb/

#ENV REDIS_URL=redis://h:password@${redisHostName}:6379
#ENV CLEARDB_DATABASE_URL=mysql://hiromaily:12345678@mysql-server/hiromaily?reconnect=true

WORKDIR /go/src/github.com/hiromaily/go-gin-wrapper

#RUN go get -d -v ./... && go build -v -o /go/bin/ginserver ./cmd/ginserver/

EXPOSE 9999

#ENTRYPOINT ["ginserver", "-f", "./configs/docker.toml"]