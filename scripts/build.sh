#!/bin/bash

. ./config/env.sh

# mkdir tmp
# chown $USER:$USER tmp
# mkdir tmp/pgdata
# chown $USER:$USER tmp/pgdata
docker-compose -f ./build/go/builder.yml build