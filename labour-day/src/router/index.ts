import {createRouter, createWebHistory, createWebHashHistory, RouteRecord, RouteRecordName} from 'vue-router'
import { setupRouterGuard } from './guard'
import {basicRoutes} from './routes'


const isHash = !!import.meta.env.VITE_APP_USE_HASH
export const router = createRouter({
  history: isHash ? createWebHashHistory('/') : createWebHistory('/'),
  routes: basicRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})


export function resetRouter() {
  router.getRoutes().forEach((route: RouteRecord) => {
    const { name } = route
    router.hasRoute(<RouteRecordName>name) && router.removeRoute(<string | symbol>name)
  })
  basicRoutes.forEach((route) => {
    !router.hasRoute(route.name) && router.addRoute(route)
  })
}

export function setupRouter() {
  setupRouterGuard(router)
  return router
}
