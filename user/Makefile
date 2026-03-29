COMMIT?=$(shell git rev-parse --short HEAD)

.PHONY: openapi-gen
openapi-gen:
	@echo "-- buliding OpenAPI"
	go generate ./cmd/...

.PHONY: build
build:
	@echo "-- building binary"
	go build \
		-o ./bin \
		./cmd