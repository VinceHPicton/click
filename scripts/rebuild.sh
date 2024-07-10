#!/bin/bash

set -e

source .env

./scripts/clean.sh
./scripts/build.sh
./scripts/up.sh