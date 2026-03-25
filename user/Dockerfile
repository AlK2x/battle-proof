FROM golang:1.26 as dependencies
WORKDIR /src/user
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /src/user
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o ./bin/user \
    ./cmd

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build --chmod=755 /src/user/bin/user /app/user
RUN ls -al .
EXPOSE 8080
ENTRYPOINT ["/app/user"]