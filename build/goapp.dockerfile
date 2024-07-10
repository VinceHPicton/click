ARG GOAPP_BASE_IMAGE

FROM ${GOAPP_BASE_IMAGE} as builder

# FROM golang:1.20 as builder

WORKDIR /opt

COPY ./cmd ./cmd
COPY ./internal ./internal

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