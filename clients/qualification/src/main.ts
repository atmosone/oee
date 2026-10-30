import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createApp } from 'vue'
import App from './App.vue'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { md3 } from 'vuetify/blueprints'

const vuetify = createVuetify({
    components,
    directives,
    blueprint: md3,
    theme: {defaultTheme: 'dark'}
})

const app = createApp(App)
app.use(vuetify)
app.mount('#app')

