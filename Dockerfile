# --- builder ---
FROM golang:1.22 as builder

ARG BUILD_VERSION
ARG BUILD_LAST_COMMIT

RUN mkdir /app
RUN mkdir /app/api-server

WORKDIR /app/api-server

COPY go.mod /app/api-server/go.mod
COPY go.sum /app/api-server/go.sum

RUN go env
RUN go env -w CGO_ENABLED=0
RUN go mod download

COPY . /app/api-server

RUN go build \
    -a \
    -o "release/api-server" \
    -ldflags " -X './internal/meta.Version=$BUILD_VERSION' -X './internal/meta.LastCommit=$BUILD_LAST_COMMIT'" \
    "./cmd/api-server"

# --- runner ---
FROM debian as runner

RUN apt update && apt upgrade -y && apt install -y ca-certificates curl && update-ca-certificates

COPY --from=builder /app/api-server/release/api-server /app/api-server/release/api-server

RUN mkdir -p /usr/local/bin
RUN ln -s /app/api-server/release/api-server /usr/local/bin/factorio-rcon-api

ENV LOG_FILE_PATH /var/log/api-server-services/api-server.log

EXPOSE 24181
EXPOSE 24180

CMD [ "/usr/local/bin/factorio-rcon-api" ]
