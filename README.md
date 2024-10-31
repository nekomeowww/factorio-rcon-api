# Factorio RCON API

> [Live API Docs](https://factorio-rcon-api.ayaka.io/apis/docs/v2)

🏭 Fully implemented wrapper for Factorio headless server console as RESTful and gRPC for easier management through APIs

## Compatible matrix

| Factorio RCON API | Factorio Server |
|-------------------|-----------------|
| 1.1.0 (`/api/v1`) | 1.x             |
| 2.0.0 (`/api/v2`) | 2.x             |

## Features

- ✅ 100% type safe, no extra type conversion needed
- ↔️ Out of the box RESTful and gRPC support
- 🎺 Native RCON protocol
- 📖 Fully API Documented

## Supported Commands

  - [x] Raw command
  - [x] Message
  - [ ] `/alerts`
  - [ ] `/enable-research-queue`
  - [ ] `/mute-programmable-speaker`
  - [ ] `/perf-avg-frames`
  - [ ] `/permissions`
  - [ ] `/reset-tips`
  - [x] `/evolution`
  - [x] `/seed`
  - [x] `/time`
  - [ ] `/toggle-action-logging`
  - [ ] `/toggle-heavy-mode`
  - [ ] `/unlock-shortcut-bar`
  - [ ] `/unlock-tips`
  - [x] `/version`
  - [x] `/admins`
  - [x] `/ban`
  - [x] `/bans`
  - [ ] `/config`
  - [ ] `/delete-blueprint-library`
  - [x] `/demote`
  - [x] `/ignore`
  - [x] `/kick`
  - [x] `/mute`
  - [x] `/mutes`
  - [x] `/promote`
  - [x] `/purge`
  - [x] `/server-save`
  - [x] `/unban`
  - [x] `/unignore`
  - [x] `/unmute`
  - [x] `/whisper`
  - [x] `/whitelist`
  - [ ] `/cheat`
  - [x] `/command` / `/c`
  - [ ] `/measured-command`
  - [ ] `/silent-command`

## Usage

> [!CAUTION]
> **Before you proceed - Security concerns**
>
> This API implementation will allow any of the users that can ACCESS the endpoint to control over & perform admin operations to Factorio server it connected to, while API server doesn't come out with any security features (e.g. Basic Auth or Authorization header based authentication).
>
> You are responsible for securing your Factorio server and the API server by either:
>
> - use Nginx/Caddy or similar servers for authentication
> - use Cloudflare Tunnel or TailScale for secure tunneling
> - use this project only for internal communication (e.g. Bots, API wrappers, admin UI, etc.)
>
> Otherwise, we are not responsible for any data loss, security breaches, save corruptions, or any other issues caused by the outside attackers.

### Pull the image

```shell
docker pull ghcr.io/nekomeowww/factorio-rcon-api
```

### Setup Factorio servers

> [!NOTE]
> **About RCON**
>
> [RCON](https://wiki.vg/RCON) is a TCP/IP-based protocol that allows server administrators to remotely execute commands, developed by Valve for Source Engine. It is widely used in game servers, including Factorio, Minecraft.

> [!CAUTION]
> **Before you proceed - Security concerns**
>
> Since RCON protocol will give administrators access to the server console, it is recommended to:
>
> - do not expose the RCON port to the public internet
> - use RCON with password authentication
> - rotate the password once a month to prevent attackers from accessing the server

When bootstraping the server, you need to specify the RCON port and password for the server to listen to with

- `--rcon-port` for the port number
- `--rcon-password` for the password

> Documentation of these parameters can be found at [Command line parameters - Factorio Wiki](https://wiki.factorio.com/Command_line_parameters)

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

### Call the API

That's it, you can now call the API with the following command:

```shell
curl -X GET http://localhost:24180/api/v2/factorio/console/command/version
```

to get the version of the Factorio game server.

## API

For API documentation, we offer Scalar powered OpenAPI UI under `/apis/docs` endpoint.

With the demo server at [https://factorio-rcon-api.ayaka.io/apis/docs/v2](https://factorio-rcon-api.ayaka.io/apis/docs/v2) live, you can view the full API documentations there, or you can run the API server locally and access the documentation at [http://localhost:24180/apis/docs/v2](http://localhost:24180/apis/docs/v2).

Alternatively, we ship the OpenAPI v2 and v3 spec in the repository:

- OpenAPI v2 spec: [v2.swagger.json](https://github.com/nekomeowww/factorio-rcon-api/blob/main/apis/factorioapi/v2/v2.swagger.json)
- OpenAPI v3 spec: [v2.swagger.v3.yaml](https://github.com/nekomeowww/factorio-rcon-api/blob/main/apis/factorioapi/v2/v2.swagger.v3.yaml)

> [!TIP]
> Additionally, we can ship the SDKs for Lua, TypeScript and Python (widely used for mods, admin panels, bots) in the future, you are welcome to contribute to the project.

For developers working with the APIs from Factorio RCON API, you can either use the above OpenAPI specs or use Protobuf files to generate types for TypeScript, Python, Go, and many more languages' SDKs with code generators. We are not going to cover all of these in this README, but you can find more information on the internet:

- [Stainless | Generate best-in-class SDKs](https://www.stainlessapi.com/) (used by OpenAI, Cloudflare, etc.)
- [Generated SDKs](https://buf.build/docs/bsr/generated-sdks/overview/)
- [Hey API](https://heyapi.dev/)

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=nekomeowww/factorio-rcon-api&type=Date)](https://star-history.com/#nekomeowww/factorio-rcon-api&Date)

## Contributors

Thanks to all the contributors!

[![contributors](https://contrib.rocks/image?repo=nekomeowww/factorio-rcon-api)](https://github.com/nekomeowww/factorio-rcon-api/graphs/contributors)
