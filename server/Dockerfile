FROM golang:alpine3.20 AS build

RUN apk add --no-cache git

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM golang:alpine3.20

RUN apk add --no-cache sqlite bash gcc build-base

COPY --from=build /go/bin/goose /usr/local/bin/goose
