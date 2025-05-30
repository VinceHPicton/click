#!/bin/bash

source .env

docker-compose -f ./docker-compose.yml down --volumes
docker volume prune -f
docker rmi -f ${PSQL_CLI_IMAGE_NAME}
docker rmi -f $(docker images --filter "dangling=true" -q)
docker rmi -f ${GOAPP_IMAGE_NAME}
docker rmi -f ${DB_CUSTOM_IMAGE_NAME}

# chmod -R +w ./tmp
# chown $USER:$USER tmp/pgdata
rm -rf ./tmp