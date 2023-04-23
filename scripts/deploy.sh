#!/bin/bash

. ./scripts/myapp.sh

export NETWORK=click_network

COMMAND="$1"
SUBCOMMAND="$2"

function localcompose(){
    local cmd=$1
    case $cmd in
        "clean")
            docker-compose -f ./deployment/compose/docker-compose.yml down
             ./scripts/psql.sh image clean
             buildMyApp clean
            rm -rf ./tmp
            ;;
        "start")
            ./scripts/psql.sh image build
            buildMyApp build
            docker-compose -f ./deployment/compose/docker-compose.yml up
            ;;
        "stop")
            docker-compose -f ./deployment/compose/docker-compose.yml down
            ;;
        *)
            echo "compose [ clean | start | stop ]"
            ;;
    esac
}

case $COMMAND in
    "local")
        localcompose $SUBCOMMAND
        ;;
    *)
        echo "$0 <command>

commands:
    compose   to spin up local network
"
        ;;
esac