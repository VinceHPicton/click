ARG PSQL_CLI_BASE_IMAGE=ubuntu:22.04

FROM ${PSQL_CLI_BASE_IMAGE}

RUN apt-get update && \
    apt-get install -y libc6-dev postgresql-client

