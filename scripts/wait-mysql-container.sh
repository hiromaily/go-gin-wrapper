#!/bin/sh

while :
do
    MSG=`docker-compose logs mysql | grep 'mysqld: ready for connections'`
    if [ -n "$MSG" ]; then
        echo 'MySQL is running!'
        sleep 1s
        #run something
        break
    else
        echo 'waiting...'
        sleep 1s
    fi

done
