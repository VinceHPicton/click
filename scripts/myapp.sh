#!/bin/bash

export GO_VER=golang:1.20
export MYAPP_IMAGE=click/myapp:current

COMMAND="$1"

function buildMyApp(){
    local cmd=$1
    case $cmd in
        "build")
        docker-compose -f ./build/go/builder.yml build
            ;;
        "clean")
            docker rmi -f ${MYAPP_IMAGE}
            ;;
        *)
            echo "[ clean | build ]"
            ;;
    esac
}
