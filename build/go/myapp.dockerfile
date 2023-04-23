ARG GO_VER

FROM ${GO_VER} as builder

ARG APP_NAME

WORKDIR /opt

COPY ./cmd ./cmd
COPY ./internal ./internal

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download && \
    env CGO_ENABLED=0 env GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./build/${APP_NAME} ./cmd/myapp

# Pack linux artefact into scratch container
FROM scratch

ARG APP_NAME

# Replace app name {goreact} here with name of your choice
COPY --from=builder /opt/build/${APP_NAME} /usr/bin/${APP_NAME}