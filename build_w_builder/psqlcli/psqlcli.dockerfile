ARG PSQL_CLI_BASE_IMAGE

FROM ${PSQL_CLI_BASE_IMAGE}

RUN apt-get update && \
    apt-get install -y libc6-dev postgresql-client

