import '@/styles/index.scss'
import 'uno.css'

import { createApp } from 'vue'
import App from './App.vue'
import { setupStore } from '@/store'
import { setupRouter } from "@/router"

const app = createApp(App)

app.use(setupStore()).use(setupRouter()).mount('#app')
