OUTPUT_PATH := /tmp/
BINARY_NAME := fourrabbitsstudio
SERVER_PORT := 8081

.PHONY: server/build
server/build:
	@go generate
	@go build -ldflags="-X main.port=${SERVER_PORT}" -o="${OUTPUT_PATH}${BINARY_NAME}" .

.PHONY: server/run
server/run: server/build 
	go run github.com/cosmtrek/air@latest
	${OUTPUT_PATH}${BINARY_NAME}

.PHONY: server/live
server/live:
	@go generate
	@podman build \
		-t ${BINARY_NAME} .
	@podman run -it --rm \
		-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
		-e AWS_REGION=${AWS_REGION} \
		-e BUCKET_NAME=${BUCKET_NAME} \
		-e MAILERLITE_TOKEN=${MAILERLITE_TOKEN} \
		-e POSTGRES_USERNAME=${POSTGRES_USERNAME} \
		-e POSTGRES_SECRET=${POSTGRES_SECRET} \
		-e POSTGRES_HOST=${POSTGRES_HOST} \
		-e POSTGRES_PORT=${POSTGRES_PORT} \
		-e POSTGRES_DB=${POSTGRES_DB} \
		-p ${SERVER_PORT}:8080 ${BINARY_NAME}


# TODO: wait for 1.22 release
# .PHONY: live
# live:
# 	@podman run -it --rm \
#     -w "/fourrabbitsstudio" \
#     -v $(shell pwd):"/fourrabbitsstudio" \
#     -p ${SERVER_PORT}:8080 \
#     cosmtrek/air

.PHONY: db/up
db/up: 
	@migrate -database cockroachdb://${POSTGRES_USERNAME}:${POSTGRES_SECRET}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=verify-full -path migrations up

.PHONY: db/down
db/down: 
	@migrate -database cockroachdb://${POSTGRES_USERNAME}:${POSTGRES_SECRET}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=verify-full -path migrations down
	
