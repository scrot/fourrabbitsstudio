OUTPUT_PATH := /tmp/
BINARY_NAME := fourrabbitsstudio
SERVER_PORT := 8081

.PHONY: build
build:
	@go build -ldflags="-X main.port=${SERVER_PORT}" -o="${OUTPUT_PATH}${BINARY_NAME}" .

.PHONY: live
live:
	@go run github.com/cosmtrek/air@latest \
		--build.cmd "make build" \
		--build.bin "${OUTPUT_PATH}${BINARY_NAME}" \
		--build.exclude_dir "" \
		--build.include_ext "go, css, tmpl"
		--build.delay "100" \
		--misc.clean_on_exit "true"


