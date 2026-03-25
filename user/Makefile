COMMIT?=$(shell git rev-parse --short HEAD)

.PHONY: build
build:
	@echo "-- building binary"
	go build \
		-o ./bin \
		./cmd