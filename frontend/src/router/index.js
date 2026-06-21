import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { title: '仪表盘' }
  },
  {
    path: '/lighting',
    name: 'Lighting',
    component: () => import('@/views/Lighting.vue'),
    meta: { title: '光照排程' }
  },
  {
    path: '/oxygen',
    name: 'Oxygen',
    component: () => import('@/views/Oxygen.vue'),
    meta: { title: '溶氧量控制' }
  },
  {
    path: '/device',
    name: 'Device',
    component: () => import('@/views/Device.vue'),
    meta: { title: '设备控制' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = `${to.meta.title} - 水族箱智能控制系统`
  next()
})

export default router
