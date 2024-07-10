#!/bin/bash

set -e

source .env

sqlc generate

docker build . -f ./build/goapp.dockerfile --build-arg GOAPP_BASE_IMAGE=${GOAPP_BASE_IMAGE} -t ${GOAPP_IMAGE_NAME}
docker build . -f ./build/dbserver.dockerfile --build-arg DB_BASE_IMAGE=${DB_BASE_IMAGE} -t ${DB_CUSTOM_IMAGE_NAME}
docker build . -f ./build/pgadmin.dockerfile --build-arg PGADMIN_BASE_IMAGE=${PGADMIN_BASE_IMAGE} -t ${PGADMIN_CUSTOM_IMAGE_NAME}



# docker-compose -f ./build/dbserver/builder.yml build
# docker-compose -f ./build/pgadmin/builder.yml build