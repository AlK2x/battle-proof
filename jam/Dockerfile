FROM golang:1.26 AS dependencies
WORKDIR /src/jam
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /src/jam
COPY . .
RUN make build
RUN make openapi-gen
RUN make test

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build --chmod=755 /src/jam/bin/jam /app/jam
RUN ls -al .
EXPOSE 8080
ENTRYPOINT ["/app/jam"]