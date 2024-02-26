OUTPUT_PATH := /tmp/
BINARY_NAME := fourrabbitsstudio
HOST := localhost
PORT := 8081

POSTGRES_DSN := ${POSTGRES_USERNAME}:${POSTGRES_SECRET}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=verify-full

.PHONY: server/build
server/build:
	@go generate ./web
	@go build -o="${OUTPUT_PATH}${BINARY_NAME}" ./cmd/server/main.go

.PHONY: server/run
server/run: server/build 
	go run github.com/cosmtrek/air@latest -c .air.toml

AIR_WORKDIR := /fourrabbitsstudio 

.PHONY: server/live
server/live:
	@podman build -f Dockerfile.dev . -t ${BINARY_NAME}-dev
	@podman run -it --rm \
		-e HOST=${HOST} \
		-e PORT=${PORT} \
		-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
		-e AWS_REGION=${AWS_REGION} \
		-e GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} \
		-e BUCKET_NAME=${BUCKET_NAME} \
		-e MAILERLITE_TOKEN=${MAILERLITE_TOKEN} \
		-e POSTGRES_USERNAME=${POSTGRES_USERNAME} \
		-e POSTGRES_SECRET=${POSTGRES_SECRET} \
		-e POSTGRES_HOST=${POSTGRES_HOST} \
		-e POSTGRES_PORT=${POSTGRES_PORT} \
		-e POSTGRES_DB=${POSTGRES_DB} \
		-e air_wd=${AIR_WORKDIR} \
		-w ${AIR_WORKDIR} \
		-v $(shell pwd):${AIR_WORKDIR} \
		-p ${SERVER_PORT}:8080 \
		${BINARY_NAME}-dev \
		-c .air.toml
		
.PHONY: db/up
db/up: 
	@migrate -database cockroachdb://${POSTGRES_DSN} -path migrations up

.PHONY: db/reset
db/reset: 
	@cockroach sql --url postgresql://${POSTGRES_DSN} --execute="TRUNCATE products, sessions, users;"
	@migrate -database cockroachdb://${POSTGRES_DSN} -path migrations down
