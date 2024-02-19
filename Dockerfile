ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app

COPY vendor vendor
COPY assets assets
COPY templates templates
COPY go.mod go.sum ./
COPY *.go .

RUN go build -ldflags="-X main.port=8080" -o=/tmp/fourrabbitsstudio .


FROM  scratch

# Import the user and group files
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Import certificates (for outgoing HTTPS requests)
COPY --from=builder /etc/ssl /etc/ssl

COPY --from=builder /tmp/fourrabbitsstudio /usr/local/bin/
CMD ["fourrabbitsstudio"]
