import {Router, RouteLocationNormalized} from "vue-router";

const baseTitle = import.meta.env.VITE_APP_TITLE

export function createPageTitleGuard(router: Router) {
  router.afterEach((to: RouteLocationNormalized) => {
    const pageTitle = to.meta?.title
    if (pageTitle) {
      document.title = `${pageTitle} | ${baseTitle}`
    } else {
      document.title = baseTitle
    }
  })
}
