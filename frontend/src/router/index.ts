import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue')
  },
  {
    path: '/products',
    name: 'ProductSelect',
    component: () => import('@/views/ProductSelect.vue')
  },
  {
    path: '/pricing',
    name: 'Pricing',
    component: () => import('@/views/Pricing.vue')
  },
  {
    path: '/competitors',
    name: 'Competitor',
    component: () => import('@/views/Competitor.vue')
  },
  {
    path: '/supply-chain',
    name: 'SupplyChain',
    component: () => import('@/views/SupplyChain.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
