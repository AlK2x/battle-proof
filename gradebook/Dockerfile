FROM golang:1.26 AS dependencies
WORKDIR /src/gradebook
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /src/gradebook
COPY . .
RUN make build
RUN make openapi-gen
RUN make test

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build --chmod=755 /src/gradebook/bin/gradebook /app/gradebook
RUN ls -al .
EXPOSE 8080
ENTRYPOINT ["/app/gradebook"]