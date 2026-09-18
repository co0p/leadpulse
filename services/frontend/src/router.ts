import { createRouter, createWebHistory } from 'vue-router'
import Home from './views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/members',
    name: 'Members',
    component: () => import('./views/Members.vue')
  },
  {
    path: '/members/add',
    name: 'AddMember',
    component: () => import('./views/AddMemberView.vue')
  },
  {
    path: '/members/:id/edit',
    name: 'EditMember',
    component: () => import('./views/EditMemberView.vue')
  },
  {
    path: '/alerts',
    name: 'Alerts',
    component: () => import('./views/Alerts.vue')
  },
  {
    path: '/reports',
    name: 'Reports',
    component: () => import('./views/Reports.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
