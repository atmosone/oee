import { defineStore } from "pinia";
import type { paths } from "./api";
import createClient from 'openapi-fetch'

const useAppStore = defineStore('app', {
    state: () => ({
        api: createClient<paths>({baseUrl: "/api"}),
        sse: new EventSource("/api/events")
    })
})

export default useAppStore;