import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../components/Home.vue'


// TODO Need to fix scroll behaviour
const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/glx/:race',
    name: 'Race',
    component: () => import(/* webpackChunkName: "choose" */ '../components/Choose.vue')
  },
  {
    path: '/play/:glxid',
    name: 'Play',
    component: () => import(/* webpackChunkName: "play" */ '../components/Play.vue')
  },
  {
    path: '/story/:glxid',
    name: 'Story',
    component: () => import(/* webpackChunkName: "play" */ '../components/Story.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
