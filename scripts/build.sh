#!/bin/bash

. ./config/env.sh

mkdir tmp
mkdir tmp/pgdata
docker-compose -f ./build/go/builder.yml build