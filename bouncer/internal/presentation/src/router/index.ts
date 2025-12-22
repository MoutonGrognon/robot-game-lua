import { createRouter, createWebHashHistory, type RouteRecordRaw } from 'vue-router'
import HomePage from '../views/HomePage.vue'

export const mainRoutes: RouteRecordRaw[] = [
  {
    path: '/home',
    name: 'Home',
    component: HomePage,
  },
  {
    path: '/ranking',
    name: 'Ranking',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/RankingPage.vue'),
  },
  {
    path: '/matches',
    name: 'Matches',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/MatchesPage.vue'),
  },
  {
    path: '/rules',
    name: 'Rules',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/RulesPage.vue'),
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/AboutUs.vue'),
  },
]

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/home',
  },
  ...mainRoutes,
  {
    path: '/matches/:id',
    name: 'Match details',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/MatchDetails.vue'),
  },
]

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
