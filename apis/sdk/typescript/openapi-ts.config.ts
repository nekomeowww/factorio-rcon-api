import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  client: '@hey-api/client-fetch',
  input: '../../factorioapi/v1/v1.swagger.v3.json',
  output: 'src/client',
})
