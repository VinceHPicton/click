#!/bin/bash

COMMAND=$1

function db(){
   createdb -h defaultserver -U postgres -p 5432 -W --maintenance-db=postgres testdb
}

function schema(){
    psql  -h defaultserver -U postgres -p 5432 -W -d testdb -f ./sql/create.sql
}

case $COMMAND in
    "db")
        db
        ;;
    "schema")
        schema
        ;;
    *)
        echo "
$0 <commands>

commands:
    db       create db
    schema   create tables
"
        ;;
esac

