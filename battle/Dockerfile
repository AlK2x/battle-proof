FROM golang:1.26 AS dependencies
WORKDIR /src/battle
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /src/battle
COPY . .
RUN make build
RUN make openapi-gen
RUN make test

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /src/battle/database/mysql /app/database/mysql
COPY --from=build --chmod=755 /src/battle/bin/battle /app/battle
RUN ls -al /app/database/mysql && \
    ls -al /app
RUN ls -al .
EXPOSE 8082
ENTRYPOINT ["/app/battle"]