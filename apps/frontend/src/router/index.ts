import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Home from '@/pages/Home.vue'
import Register from '@/pages/Register.vue'
import Login from '@/pages/Login.vue'
import RegisterConfirm from '@/pages/RegisterConfirm.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  { 
    path: '/login',
    name: 'Login', 
    component: Login
  },
  { 
    path: '/register',
    name: 'Register', 
    component: Register 
  },
  { 
    path: '/register/confirm', 
    name: 'RegisterConfirm', 
    component: RegisterConfirm 
  }
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})