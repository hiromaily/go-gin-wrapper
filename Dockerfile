# Dcokerfile for go-gin-wrapper
FROM golang:1.15.6-alpine3.12

RUN apk add --no-cache git mysql-client

#ARG redisHostName=default-redis-server
#ARG mysqlHostName=default-mysql-server
RUN mkdir -p /go/src/github.com/hiromaily/go-gin-wrapper/tmp/ /var/log/goweb/

#ENV REDIS_URL=redis://h:password@${redisHostName}:6379
#ENV CLEARDB_DATABASE_URL=mysql://hiromaily:12345678@mysql-server/hiromaily?reconnect=true

WORKDIR /go/src/github.com/hiromaily/go-gin-wrapper
COPY . .

RUN GO111MODULE=off go get -u github.com/oxequa/realize
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/ginserver ./cmd/ginserver/

EXPOSE 8080
CMD ["/go/bin/ginserver", "-f", "./configs/docker.toml", "-crypto"]
