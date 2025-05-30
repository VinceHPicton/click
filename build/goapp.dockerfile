ARG GOAPP_BASE_IMAGE

FROM ${GOAPP_BASE_IMAGE}-alpine as builder

# FROM golang:1.20 as builder

WORKDIR /opt

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./db ./db
COPY sqlc.yaml .

# Install dependencies
RUN apk add --no-cache curl 

# ENV SQLC_VERSION=1.29.0
ARG SQLC_VERSION

# Download and install sqlc
RUN curl -L https://github.com/sqlc-dev/sqlc/releases/download/v${SQLC_VERSION}/sqlc_${SQLC_VERSION}_linux_amd64.tar.gz \
    | tar -xz -C /usr/local/bin

RUN sqlc generate

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download && \
    env CGO_ENABLED=0 env GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./build/goapp ./cmd/goapp

# Pack linux artefact into scratch container
# FROM scratch
FROM alpine

ARG APP_NAME

COPY --from=builder /opt/build/goapp /usr/bin/goapp

ENTRYPOINT ./usr/bin/goapp