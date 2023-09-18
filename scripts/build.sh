#!/bin/bash

. ./config/env.sh

docker-compose -f ./build/go/builder.yml build