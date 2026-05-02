FROM golang:1.26 AS dependencies
WORKDIR /src/user
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /src/user
COPY . .
RUN make build
RUN make openapi-gen
RUN make test

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /src/user/database/mysql /app/database/mysql
COPY --from=build --chmod=755 /src/user/bin/user /app/user
RUN ls -al /app/database/mysql && \
    ls -al /app
RUN ls -al .
EXPOSE 8080
ENTRYPOINT ["/app/user"]