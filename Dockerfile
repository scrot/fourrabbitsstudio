ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app

COPY vendor vendor
COPY assets assets
COPY templates templates
COPY go.mod go.sum ./
COPY *.go .

# TODO: add npm buildstage for go generate

RUN go build -ldflags="-X main.port=8080" -o=/tmp/fourrabbitsstudio .

FROM  gcr.io/distroless/base:latest

COPY --from=builder /tmp/fourrabbitsstudio /usr/local/bin/

CMD ["fourrabbitsstudio"]
