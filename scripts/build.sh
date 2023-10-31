#!/bin/bash

set -e

. ./config/env.sh

sqlc generate
docker-compose -f ./build/go/builder.yml build