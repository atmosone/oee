import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createRouter, createWebHistory } from 'vue-router'

const app = createApp(App)

const pinia = createPinia()

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'start',
      component: () => import('@/AppStartPage.vue'),
    },
    {
      path: '/equipment/:id',
      name: 'equipment',
      component: () => import('@/equipment/EquipmentPage.vue'),
    },
    {
      path: '/lines/:id',
      name: 'equipment_line',
      component: () => import('@/lines/LinePage.vue'),
    },
    {
      path: '/specifications/:id',
      name: 'specification',
      component: () => import('@/specifications/SpecificationPage.vue'),
    },
  ],
})

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light',
  },
})

app.use(pinia)
app.use(router)
app.use(vuetify)
app.mount('#app')
