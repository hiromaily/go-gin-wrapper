#!/bin/sh
sleep 10s
mongo 127.0.0.1:${MONGO_PORT}/admin --eval "var port = ${MONGO_PORT};" ./init.js
mongorestore -h 127.0.0.1:${MONGO_PORT} --db hiromaily dump/hiromaily
