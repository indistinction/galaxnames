import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import VueSocialSharing from 'vue-social-sharing'

const app = createApp(App)

app.config.globalProperties.$backendhost = "https://mini1.galaxiators.net" //"http://127.0.0.1:8001"
app.config.globalProperties.$activeRaces = ["c"]
//"abysstea", "egregore", "litagiar", "satyne", "thandrac"

app.use(VueSocialSharing);
app.use(router).mount('#app')






