#!/bin/bash

source .env

docker-compose -f ./docker-compose.yml up -d --wait 
sleep 2
./scripts/migrate.sh