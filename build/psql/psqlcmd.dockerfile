ARG OS_VER

# App container
FROM ${OS_VER}

RUN apt-get update && \
    apt-get install -y libc6-dev postgresql-client

