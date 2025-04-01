# --- builder ---
FROM golang:1.24 as builder

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
RUN apt install -y lsof net-tools iproute2 telnet procps

COPY --from=builder /app/api-server/release/api-server /app/api-server/release/api-server

RUN mkdir -p /usr/local/bin
RUN ln -s /app/api-server/release/api-server /usr/local/bin/factorio-rcon-api

WORKDIR /app/api-server
RUN mkdir -p /app/api-server/logs

EXPOSE 24180
EXPOSE 24181

CMD [ "/usr/local/bin/factorio-rcon-api" ]
