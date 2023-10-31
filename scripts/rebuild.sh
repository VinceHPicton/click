#!/bin/bash

set -e

. ./config/env.sh

./scripts/clean.sh
./scripts/build.sh
./scripts/up.sh