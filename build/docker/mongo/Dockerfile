# Dcokerfile for Mongo

FROM mongo:4.2-bionic

ARG mongo_port=27017
ENV MONGO_PORT=$mongo_port

COPY init.js /docker-entrypoint-initdb.d/
#COPY mongod.conf /etc/

EXPOSE 27017

CMD ["mongod", "--auth"]
