#!/bin/bash

. ./config/env.sh

docker-compose -f ./build/psqlcli/builder.yml build # --> this is for psql commands, not pgadmin or db server
docker-compose -f ./build/go/builder.yml build