OUTPUT_PATH := /tmp/
BINARY_NAME := fourrabbitsstudio
SERVER_PORT := 8081

.PHONY: build
build:
	@go generate
	@go build -ldflags="-X main.port=${SERVER_PORT}" -o="${OUTPUT_PATH}${BINARY_NAME}" .

.PHONY: run
run: build 
	go run github.com/cosmtrek/air@latest
	${OUTPUT_PATH}${BINARY_NAME}

.PHONY: live
live:
	@go generate
	@podman build \
		-t ${BINARY_NAME} .
	@podman run -it --rm \
		-e BUCKET_NAME=${BUCKET_NAME} \
		-e MAILERLITE_TOKEN=${MAILERLITE_TOKEN} \
		-e POSTGRES_DSN=${POSTGRES_DSN} \
		-v ${HOME}/.aws/config:/root/.aws/config:ro \
		-v ${HOME}/.aws/credentials:/root/.aws/credentials:ro \
		-p ${SERVER_PORT}:8080 ${BINARY_NAME}


# TODO: wait for 1.22 release
# .PHONY: live
# live:
# 	@podman run -it --rm \
#     -w "/fourrabbitsstudio" \
#     -v $(shell pwd):"/fourrabbitsstudio" \
#     -p ${SERVER_PORT}:8080 \
#     cosmtrek/air

