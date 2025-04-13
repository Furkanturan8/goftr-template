import type { App } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from './routes'
import { useUserStore } from '@/store/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Navigation Guard
// kullanıcı giriş yapmadıysa ve başka yere istek atıyorsa login sayfasına yönlendir
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const publicPages = ['/login', '/register']
  const authRequired = !publicPages.includes(to.path)

  if (authRequired && !userStore.isAuthenticated) {
    const redirect = to.fullPath
    const urlParams = new URLSearchParams()
    urlParams.set('redirect', redirect)
    return next(`/login?${urlParams.toString()}`)
  }

  next()
})

export default function (app: App) {
  app.use(router)
}

export { router }
