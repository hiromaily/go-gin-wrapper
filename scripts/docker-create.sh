#!/bin/sh

# This code is outdated, not used anymore

###############################################################################
# Using docker-composer for go-gin-wrapper
###############################################################################
#echo $RUN_TEST

###############################################################################
# Environment
###############################################################################
CONTAINER_NAME=web
CONTAINER2_NAME=web-redis
CONTAINER3_NAME=web-mysql
IMAGE_NAME=go-gin-wrapper:v1.1


###############################################################################
# Remove Container And Image
###############################################################################
DOCKER_PSID=`docker ps -af name="${CONTAINER_NAME}" -q`
if [ ${#DOCKER_PSID} -ne 0 ]; then
    docker rm -f ${CONTAINER_NAME}
fi

DOCKER_PSID=`docker ps -af name="${CONTAINER2_NAME}" -q`
if [ ${#DOCKER_PSID} -ne 0 ]; then
    docker rm -f ${CONTAINER2_NAME}
fi

DOCKER_PSID=`docker ps -af name="${CONTAINER3_NAME}" -q`
if [ ${#DOCKER_PSID} -ne 0 ]; then
    docker rm -f ${CONTAINER3_NAME}
fi

DOCKER_PSID=`docker ps -af name="${CONTAINER4_NAME}" -q`
if [ ${#DOCKER_PSID} -ne 0 ]; then
    docker rm -f ${CONTAINER4_NAME}
fi

#docker rm -f $(docker ps -aq)

DOCKER_IMGID=`docker images "${IMAGE_NAME}" -q`
if [ ${#DOCKER_IMGID} -ne 0 ]; then
    docker rmi ${IMAGE_NAME}
fi


###############################################################################
# Docker-compose / build and up
###############################################################################
docker-compose  build
docker-compose  up -d

if [ $RUN_TEST -eq 1 ]; then
    # test mode
    sleep 1s

    # create test data on docker container mysql
    export DB_NAME=hiromaily2
    export DB_PORT=23306
    export DB_USER=root
    export DB_PASS=root
    sh ./scripts/create-test-db.sh
    #mysql -uroot -proot -h127.0.0.1 -P13306 < ./tests/createdb.sql

    docker exec -it ${CONTAINER_NAME} /bin/bash -c "
        export RUN_TEST=1;
        go get -d -v ./...;
        go test -v cmd/ginserver/*.go -f ../../configs/docker.toml -crypto;
    "
    #docker exec -it web bash ./docker-entrypoint.sh
else
    # run server mode
    # foreground
    #docker exec -it ${CONTAINER_NAME} bash ./docker-entrypoint.sh

    # background(trying)
    docker exec -itd ${CONTAINER_NAME} bash ./docker-entrypoint.sh
fi
###############################################################################
# Docker-compose / check
###############################################################################
docker-compose ps
docker-compose logs


###############################################################################
# Exec
###############################################################################
#docker exec -it web bash


###############################################################################
# Test
###############################################################################



###############################################################################
# Docker-compose / down
###############################################################################
#docker-compose -f ${COMPOSE_FILE} down

###############################################################################
# Check connection
###############################################################################
#mysql -u root -p -h 127.0.0.1 -P 13306
#redis-cli -h 127.0.0.1 -p 16379 -a password

#Access by browser
#http://docker.hiromaily.com:9999/
