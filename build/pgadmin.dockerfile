ARG PGADMIN_BASE_IMAGE

FROM ${PGADMIN_BASE_IMAGE}

ARG PGADMIN_SERVERS_FILENAME

COPY ${PGADMIN_SERVERS_FILENAME} /pgadmin4/servers.json
