#!/bin/bash

. ./config/env.sh

# docker run --rm --net=host \
#   -v $(pwd)/migrations:/migrations \
#   arigaio/atlas migrate apply \
#   --url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

atlas migrate apply --url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" --dir "file://db/migrations"