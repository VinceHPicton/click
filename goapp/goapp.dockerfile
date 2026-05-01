ARG GOAPP_BASE_IMAGE=golang:1.26

FROM ${GOAPP_BASE_IMAGE} as builder

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /opt

COPY ./goapp/cmd ./cmd
COPY ./goapp/internal ./internal
COPY ./db ./db
COPY sqlc.yaml .

RUN sqlc generate

COPY ./goapp/go.mod ./go.mod
COPY ./goapp/go.sum ./go.sum

RUN go mod download
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./build/goapp ./cmd/main.go

# Pack linux artefact into scratch container
FROM alpine

ARG APP_NAME

COPY --from=builder /opt/build/goapp /usr/bin/goapp

ENTRYPOINT ./usr/bin/goapp