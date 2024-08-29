# `stdserve`

🎺 stdin & stdout API wrapper and protocol for Factorio & Minecraft server management

## Demo

https://github.com/user-attachments/assets/1e955e70-a2da-488f-9622-cee753624c29

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

- OpenAPI v2 spec: [v1.swagger.json](https://github.com/nekomeowww/stdserve/blob/main/apis/stdserveapi/v1/v1.swagger.json)
- OpenAPI v3 spec: [v1.swagger.v3.yaml](https://github.com/nekomeowww/stdserve/blob/main/apis/stdserveapi/v1/v1.swagger.v3.yaml)
