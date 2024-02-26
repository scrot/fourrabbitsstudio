OUTPUT_PATH := /tmp/
BINARY_NAME := fourrabbitsstudio
SERVER_PORT := 8081

POSTGRES_DSN := ${POSTGRES_USERNAME}:${POSTGRES_SECRET}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=verify-full

.PHONY: server/build
server/build:
	@go generate
	@go build -ldflags="-X main.port=${SERVER_PORT}" -o="${OUTPUT_PATH}${BINARY_NAME}" .

.PHONY: server/run
server/run: server/build 
	go run github.com/cosmtrek/air@latest
	${OUTPUT_PATH}${BINARY_NAME}

AIR_WORKDIR := /go/src/github.com/scrot/fourrabbitsstudio 

.PHONY: server/live
server/live:
	@podman build -f Dockerfile.dev . -t ${BINARY_NAME}-dev
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
		-e AIR_WD=${AIR_WORKDIR} \
		-w ${AIR_WORKDIR} \
    -v $(shell pwd):${AIR_WORKDIR} \
		-p ${SERVER_PORT}:8080 \
		${BINARY_NAME}-dev
		
.PHONY: db/up
db/up: 
	@migrate -database cockroachdb://${POSTGRES_DSN} -path migrations up

.PHONY: db/reset
db/reset: 
	@cockroach sql --url postgresql://${POSTGRES_DSN} --execute="TRUNCATE products, sessions, users;"
	@migrate -database cockroachdb://${POSTGRES_DSN} -path migrations down
