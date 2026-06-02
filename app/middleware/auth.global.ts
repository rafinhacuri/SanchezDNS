export default defineNuxtRouteMiddleware(async (to) => {
  const { refresh, isLoggedIn, user } = useUser()

  // * Atualizando a sessão
  await refresh()

  if (!isLoggedIn.value) {
    return navigateTo('/login')
  }

  if (to.path === '/login') {
    return navigateTo('/')
  }

  if (to.path === '/logs' || (to.path === '/cadastros' && user.value.level !== 'admin')) {
    return navigateTo('/')
  }
})
