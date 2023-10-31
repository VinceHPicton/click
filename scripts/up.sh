#!/bin/bash

. ./config/env.sh

docker-compose -f ./docker-compose.yml up -d --wait 
sleep 2
./scripts/migrate.sh