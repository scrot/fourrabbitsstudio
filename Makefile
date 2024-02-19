OUTPUT_PATH := /tmp
BINARY_NAME := fourrabbitsstudio
SERVER_PORT := 8081

.PHONY: build
build:
	@go build -ldflags="-X main.port=8080" -o="${OUTPUT_PATH}${BINARY_NAME}" .

.PHONY: run
run:
	@ podman build \
		-t ${BINARY_NAME} .
	@ podman run -it \
		-v ${HOME}/.aws/credentials:/root/.aws/credentials:ro \
		-p ${SERVER_PORT}:8080 ${BINARY_NAME}


# TODO: wait for 1.22 release
.PHONY: live
live:
	@podman run -it --rm \
    -w "/fourrabbitsstudio" \
    -v $(shell pwd):"/fourrabbitsstudio" \
    -p ${SERVER_PORT}:8080 \
    cosmtrek/air

