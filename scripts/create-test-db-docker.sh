#!/bin/sh

echo 'create test db on docker'
docker-compose exec mysql mysql -u root -proot  -e "$(cat ./test/sql/create_gogin-test.sql)"
docker-compose exec mysql mysql -u root -proot  -e "$(cat ./test/sql/data_gogin-test.sql)"
