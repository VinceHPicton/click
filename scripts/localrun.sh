#!/bin/bash

set -e

source .env

# Needs to be on same network as db container, best to pass flag args to the .sh running command to post the right values in, 
# The container from this lingers after use, too annoying to bother using.
# Also, you cant use the container name to give it the DB address

# This would only be useful if the DB was not running as a container

if [ $# -eq 0 ]; then
    echo "Usage: $0 <database_address>"
    exit 1
fi

DB_ADDR=$1

docker run \
 -e DB_ADDR=$1 \
 -e DB_PORT=${DB_PORT} \
 -e DB_USER=${DB_USER} \
 -e DB_PASSWORD=${DB_PASSWORD} \
 -e DB_NAME=${DB_NAME} \
 -e GOAPP_PORT=${GOAPP_PORT} \
 -p ${GOAPP_PORT}:${GOAPP_PORT} \
 --network ${NETWORK_NAME} \
 --name local-goapp \
 ${GOAPP_IMAGE_NAME}