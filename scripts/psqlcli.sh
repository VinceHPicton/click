#!/bin/bash

source .env

docker run --network=${NETWORK} -v ${PWD}/db/psql/config/.pgpass:/root/.pgpass -v ${PWD}/db/psql/sql:/opt/sql/ -v ${PWD}/db/psql/scripts:/opt/scripts -w /opt -it --rm ${PSQL_CLI_IMAGE_NAME} /bin/bash
