OUTPUT_PATH := ./tmp/
BINARY_NAME := server
SERVER_PORT := 8081

POSTGRES_DSN := ${POSTGRES_USERNAME}:${POSTGRES_SECRET}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=verify-full

.PHONY: server/build
server/build:
	@go mod tidy && go mod vendor
	@go generate ./web
	@go build -buildvcs=false -o=${OUTPUT_PATH}${BINARY_NAME} ./cmd/server/main.go

.PHONY: server/run
server/run: server/build 
	@go run github.com/cosmtrek/air@latest -c .air.toml

.PHONY: server/live
server/live:
	@podman build -f Dockerfile.dev . -t fourrabbitsstudio-dev
	@podman run -it --rm \
		-e BUCKET_NAME=${BUCKET_NAME} \
		-e GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} \
		-e MAILERLITE_TOKEN=${MAILERLITE_TOKEN} \
		-e POSTGRES_USERNAME=${POSTGRES_USERNAME} \
		-e POSTGRES_SECRET=${POSTGRES_SECRET} \
		-e POSTGRES_HOST=${POSTGRES_HOST} \
		-e POSTGRES_PORT=${POSTGRES_PORT} \
		-e POSTGRES_DB=${POSTGRES_DB} \
		-e air_wd=/src \
		-w /src \
		-v $(shell pwd):/src \
		-p ${SERVER_PORT}:8080 \
		fourrabbitsstudio-dev \
		-c .air.toml
		
.PHONY: db/up
db/up: 
	@migrate -database cockroachdb://${POSTGRES_DSN} -path migrations up

.PHONY: db/reset
db/reset: 
	@cockroach sql --url postgresql://${POSTGRES_DSN} --execute="TRUNCATE products, sessions, users;"
	@migrate -database cockroachdb://${POSTGRES_DSN} -path migrations down
