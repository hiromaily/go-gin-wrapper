# Dcokerfile for Mongo

FROM mongo:3.4

ARG mongo_port=27017
ENV MONGO_PORT=$mongo_port

RUN mkdir settings
WORKDIR /settings
COPY . .
