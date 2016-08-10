#!/bin/sh

###############################################################################
# Using docker-composer for go-gin-wrapper
###############################################################################

###############################################################################
# Environment
###############################################################################
CONTAINER_NAME=web
IMAGE_NAME=go-gin-wrapper:v1.1


###############################################################################
# Remove Container And Image
###############################################################################
#DOCKER_PSID=`docker ps -af name="${CONTAINER_NAME}" -q`
#if [ ${#DOCKER_PSID} -ne 0 ]; then
#    docker rm -f ${CONTAINER_NAME}
#fi
docker rm -f $(docker ps -aq)

DOCKER_IMGID=`docker images "${IMAGE_NAME}" -q`
if [ ${#DOCKER_IMGID} -ne 0 ]; then
    docker rmi ${IMAGE_NAME}
fi


###############################################################################
# Docker-compose / build and up
###############################################################################
docker-compose  build
docker-compose  up -d

#settings: moved to docker-compose.yml
#docker exec -it ${CONTAINER_NAME} bash ./docker-entrypoint.sh

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
