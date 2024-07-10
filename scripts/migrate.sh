#!/bin/bash

source .env

# This achieves the same as the below command but using docker
# docker run --rm --net=host \
#   -v $(pwd)/migrations:/migrations \
#   arigaio/atlas migrate apply \
#   --url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

# This command will run any migration files found in the --dir directory against the database at --url that have not already been run
atlas migrate apply \
--url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" \
--dir "file://db/migrations"