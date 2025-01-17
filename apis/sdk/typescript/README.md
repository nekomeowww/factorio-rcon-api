# `factorio-rcon-api-client`

Generated TypeScript SDK client for [https://github.com/nekomeowww/factorio-rcon-api](https://github.com/nekomeowww/factorio-rcon-api)

## Getting started

```shell
ni factorio-rcon-api-client -D # from @antfu/ni, can be installed via `npm i -g @antfu/ni`
pnpm i factorio-rcon-api-client -D
yarn i factorio-rcon-api-client -D
npm i factorio-rcon-api-client -D
```

## Usage

```ts
// Import the client and the API functions
import { client, v2FactorioConsoleCommandRawPost } from 'factorio-rcon-api-client'

async function main() {
  // Set the base URL of the API
  client.setConfig({
    baseUrl: 'http://localhost:3000',
  })

  // Call POST /api/v2/factorio/console/command/raw
  const res = await v2FactorioConsoleCommandRawPost({
    body: {
      input: '/help', // The command to run
    },
  })

  console.log(res) // The response from the API
}

main().catch(console.error)
```
