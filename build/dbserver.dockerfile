ARG DB_BASE_IMAGE

FROM ${DB_BASE_IMAGE}

COPY /db/seeds /docker-entrypoint-initdb.d/
COPY /db/schema/0_schema.sql /docker-entrypoint-initdb.d/