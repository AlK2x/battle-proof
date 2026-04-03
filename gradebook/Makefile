COMMIT?=$(shell git rev-parse --short HEAD)

.PHONY: openapi-gen
openapi-gen:
	@echo "-- buliding OpenAPI"
	go generate ./cmd/...

.PHONY: build
build:
	@echo "-- building binary"
	CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o ./bin/gradebook \
    ./cmd

.PHONY: test
test:
	@echo "-- running tests"
	go test -v -race ./...

.PHONY: docker-build
docker-build:
	@echo "-- building Docker image"
	docker build --progress=plain -t gradebook-app:latest .