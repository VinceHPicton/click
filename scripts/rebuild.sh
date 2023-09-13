#!/bin/bash

. ./config/env.sh

./scripts/clean.sh
./scripts/build.sh
./scripts/up.sh