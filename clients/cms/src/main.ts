import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createApp, type InjectionKey } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import createClient from 'openapi-fetch'
import { type paths } from './schema.ts'

const api = createClient<paths>({
    baseUrl: "/api"
})

const vuetify = createVuetify({
    components,
    directives,
    theme: {
        defaultTheme: 'light'
    }
})

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
        path: "/equipment",
        component: () => import("@/pages/EquipmentPage.vue"),
    },
    {
        path: "/lines",
        component: () => import("@/pages/LinesPage.vue"),
        children: [
            {
                path: ":line_id",
                component: () => import("@/pages/LinePage.vue"),
                props: true,
            }
        ]
    }
  ],
})

const app = createApp(App)
app.config.globalProperties.$api = api
app.use(createPinia())
app.use(router)
app.use(vuetify)
app.mount('#app')
