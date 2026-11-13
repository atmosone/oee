import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createApp } from 'vue'
import App from './App.vue'
import {createVuetify} from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'

const pinia = createPinia()

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/", 
      component: () => import('@/screens/StartPage.vue')
    },
    {
      path: "/equipment", 
      component: () => import('@/screens/EquipmentsScreen.vue'),
      children: [
        {
          path: ":id",
          component: () => import('@/screens/EquipmentScreen.vue'),
        }
      ]
    },
    {
      path: "/lines", 
      component: () => import('@/screens/LinesScreen.vue'),
      children: [
        {
          path: ":id",
          component: () => import('@/screens/LineScreen.vue'),
        }
      ]
    },
    {
      path: "/specifications", 
      component: () => import('@/screens/SpecificationsScreen.vue'),
      children: [
        {
          path: ":id",
          component: () => import('@/screens/SpecificationScreen.vue'),
        }
      ]
    },
  ],
})

const vuetify = createVuetify({
    components,
    directives,
})

const app = createApp(App)
app.use(pinia)
app.use(router)
app.use(vuetify)
app.mount('#app')
