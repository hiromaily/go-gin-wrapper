#!/bin/sh

echo 'create test db on docker'
docker-compose exec mysql mysql -u root -proot  -e "$(cat ./test/sql/create_go-gin-test.sql)"
docker-compose exec mysql mysql -u root -proot  -e "$(cat ./test/sql/data_go-gin-test.sql)"
