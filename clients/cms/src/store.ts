import type { paths } from "@/schema";
import createClient from "openapi-fetch";
import { defineStore } from "pinia";

const useAppStore = defineStore("app", {
    state: () => ({
        client: createClient<paths>({baseUrl: '/api'})
    }),
})

export default useAppStore;