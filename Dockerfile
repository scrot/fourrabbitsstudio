ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app

COPY cmd cmd
COPY web web
COPY internal internal
COPY vendor vendor
COPY go.mod go.sum ./

RUN go build -o=/tmp/fourrabbitsstudio ./cmd/server

FROM  gcr.io/distroless/base:latest

COPY --from=builder /tmp/fourrabbitsstudio /usr/local/bin/

CMD ["fourrabbitsstudio"]
