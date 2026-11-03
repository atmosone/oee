import createClient from 'openapi-fetch'
import { defineStore } from 'pinia'
import type { paths } from './schema'

const useAppStore = defineStore('app', {
  state: () => ({
    client: createClient<paths>({ baseUrl: '/api' }),
  }),
})

export default useAppStore
