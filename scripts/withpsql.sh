#!/bin/bash

source .env

./scripts/build.sh
docker build . -f ./build/psqlcli.dockerfile --build-arg PSQL_CLI_BASE_IMAGE=${PSQL_CLI_BASE_IMAGE} -t ${PSQL_CLI_IMAGE_NAME}

# docker-compose -f ./build/psqlcli/builder.yml build