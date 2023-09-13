#!/bin/bash

. ./config/env.sh

docker-compose -f ./docker-compose.yml down
docker rmi -f ${PSQL_CLI_IMAGE_NAME}
docker rmi -f $(docker images --filter "dangling=true" -q)
docker rmi -f ${GOAPP_IMAGE_NAME}
rm -rf ./tmp