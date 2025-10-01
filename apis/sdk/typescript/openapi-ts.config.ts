import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: '../../factorioapi/v1/v1.swagger.v3.json',
  output: 'src/client',
})
