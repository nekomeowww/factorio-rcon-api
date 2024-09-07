# Factorio RCON API

🏭 Fully implemented wrapper for Factorio headless server console as RESTful and gRPC for easier management through APIs

## Usage

### Pull the image

```shell
docker pull ghcr.io/nekomeowww/factorio-rcon-api
```

### Setup servers

When bootstraping the server, you need to specify the RCON port and password for the server to listen to with

- `--rcon-port` for the port number
- `--rcon-password` for the password

The command may look like this:

```shell
./factorio \
 --start-server /path/to/saves/my-save.zip \
 --rcon-port 27015 \
 --rcon-password 123456
```

Or on macOS:

```shell
~/Library/Application\ Support/Steam/steamapps/common/Factorio/factorio.app/Contents/MacOS/factorio \
 --start-server /path/to/saves/my-save.zip \
 --rcon-port 27015 \
 --rcon-password 123456
```

Once you are ready, go to the next step to start the API server.

### Factorio server ran with Docker

This is kind of hard to make them communicate easily.

We will need to create dedicated network for the containers to communicate with each other.

```shell
docker network create factorio
```

Then, obtain the IP address of the Factorio server container.

```shell
docker container inspect factorio-server --format '{{ .NetworkSettings.Networks.factorio.IPAddress }}'
```

Then, start the API server with the following command with the IP address obtained:

```shell
docker run \
  --rm \
  -e FACTORIO_RCON_HOST=<factorio-server-ip> \
  -e FACTORIO_RCON_PORT=27015 \
  -e FACTORIO_RCON_PASSWORD=123456 \
  -p 24180:24180 \
  ghcr.io/nekomeowww/factorio-rcon-api:unstable
```

### Factorio server not ran with Docker, Factorio RCON API ran with Docker

For running Factorio server and Factorio RCON API in a same server while not having Factorio server in Docker, you can start the API server with the following command:

```shell
docker run \
  --rm \
  -e FACTORIO_RCON_HOST=host.docker.internal \
  -e FACTORIO_RCON_PORT=27015 \
  -e FACTORIO_RCON_PASSWORD=123456 \
  -p 24180:24180 \
  ghcr.io/nekomeowww/factorio-rcon-api:unstable
```

## API

- OpenAPI v2 spec: [v1.swagger.json](https://github.com/nekomeowww/stdserve/blob/main/apis/stdserveapi/v1/v1.swagger.json)
- OpenAPI v3 spec: [v1.swagger.v3.yaml](https://github.com/nekomeowww/stdserve/blob/main/apis/stdserveapi/v1/v1.swagger.v3.yaml)
