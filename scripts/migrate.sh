#!/bin/bash

. ./config/env.sh

docker run --rm --net=host \
  -v $(pwd)/db/migrations:/migrations \
  arigaio/atlas migrate apply
  --url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"