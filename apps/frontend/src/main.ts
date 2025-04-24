import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
import { router } from './router'
import 'aos/dist/aos.css'
import AOS from 'aos'

// Tailwind CSS
import './styles/components.css'

// AOS (Animate On Scroll)の導入
AOS.init()

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.mount('#app')