# --- builder ---
FROM golang:1.23@sha256:73f06be4578c9987ce560087e2e2ea6485fb605e3910542cadd8fa09fc5f3e31 as builder

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
FROM debian@sha256:10901ccd8d249047f9761845b4594f121edef079cfd8224edebd9ea726f0a7f6 as runner

RUN apt update && apt upgrade -y && apt install -y ca-certificates curl && update-ca-certificates

COPY --from=builder /app/api-server/release/api-server /app/api-server/release/api-server

RUN mkdir -p /usr/local/bin
RUN ln -s /app/api-server/release/api-server /usr/local/bin/factorio-rcon-api

WORKDIR /app/api-server
RUN mkdir -p /app/api-server/logs

EXPOSE 24180
EXPOSE 24181

CMD [ "/usr/local/bin/factorio-rcon-api" ]
