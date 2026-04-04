COMMIT?=$(shell git rev-parse --short HEAD)
GOLANGCI_LINT_VERSION ?= v2.11.4

.PHONY: openapi-gen
openapi-gen:
	@echo "-- buliding OpenAPI"
	go generate ./cmd/...

.PHONY: build
build:
	@echo "-- building binary"
	CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o ./bin/user \
    ./cmd

.PHONY: test
test:
	@echo "-- running tests"
	go test -v -race ./...

.PHONY: docker-build
docker-build:
	@echo "-- building Docker image"
	docker build --progress=plain -t user-app:latest .

.PHONY: lint
lint:
	@echo "-- lint code"
	docker run --rm -t \
	  -v "$(PWD):/app" -w /app \
	  golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
	  golangci-lint run ./...