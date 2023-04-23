#!/bin/bash

export PSQL_CLI_IMAGE=click/psqlcmd:current
export PSQL_CLI_CONTAINER=psqlcli
export NETWORK=click_network

COMMAND="$1"
SUBCOMMAND="$2"

function image(){
    local cmd="$1"
    case $cmd in
        "build")
            docker-compose -f ./build/psql/builder.yml build
            ;;
        "clean")
            docker rmi -f ${PSQL_CLI_IMAGE}
            docker rmi -f $(docker images --filter "dangling=true" -q)
            ;;
        *)
            echo "image [ build | clean]"
            ;;
    esac
}

function client(){
    local cmd=$1
    case $cmd in
        "cli")
            docker run --network=${NETWORK} -v ${PWD}/db/psql/config/.pgpass:/root/.pgpass -v ${PWD}/db/psql/sql:/opt/sql/ -v ${PWD}/db/psql/scripts:/opt/scripts -w /opt -it --rm ${PSQL_CLI_IMAGE} /bin/bash
            ;;
        *)
            echo "client [ cli ]"
            ;;
    esac
}

case $COMMAND in
    "client")
        client $SUBCOMMAND
        ;;
    "clean")
        image clean
        ;;
    "image")
        image $SUBCOMMAND
        ;;
    *)
        echo "$0 <commands>

commands:  
    clean    to reset psql deployment to empty state
    client   spins up a client app to interact with psql server
    image    build and clean docker images
"
        ;;
esac