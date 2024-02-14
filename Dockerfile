ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine as builder
WORKDIR /usr/src/app
COPY . .
RUN go build -v -o /run-app .


FROM scratch
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
