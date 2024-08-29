# `stdserve`

🎺 stdin & stdout API wrapper and protocol for Factorio & Minecraft server management

## Usage

### Factorio server

```shell
go run ./cmd/stdserve serve -p factorio -- '~/Library/Application Support/Steam/steamapps/common/Factorio/factorio.app/Contents/MacOS/factorio' --start-server './saves/my-save.zip'
```

## API

### `POST /api/v1/stdin/execute`

```shell
curl -X POST -H "Content-Type: application/json" -d '{"input": "/help"}' http://localhost:10080/api/v1/stdin/execute
```
