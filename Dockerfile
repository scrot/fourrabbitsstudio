ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app

COPY vendor vendor
COPY assets assets
COPY templates templates
COPY go.mod go.sum ./
COPY *.go .

RUN go build -ldflags="-X main.port=8080" -o=/tmp/fourrabbitsstudio .

FROM  gcr.io/distroless/static:latest

COPY --from=builder /tmp/fourrabbitsstudio /usr/local/bin/

# USER nonroot

CMD ["fourrabbitsstudio"]
