#!/bin/bash

export NETWORK=click_network

export GOAPP_BASE_IMAGE=golang:1.20 # PREVIOUSLY GO_VER, adjust for new Go versions when needed
export GOAPP_IMAGE_NAME=click/goapp:current

export PSQL_CLI_BASE_IMAGE=ubuntu:22.04 # Previously OS_VER
export PSQL_CLI_IMAGE_NAME=click/psqlcli:current
export PSQL_CLI_CONTAINER=psqlcli
