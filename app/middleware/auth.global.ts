export default defineNuxtRouteMiddleware(async (to) => {
  const { refresh, isLoggedIn, user } = useUser()

  // * Atualizando a sessão
  await refresh()

  if (!isLoggedIn.value && to.path !== '/login') {
    return navigateTo('/login')
  }

  if (isLoggedIn.value && to.path === '/login') {
    return navigateTo('/')
  }

  if (
    (to.path === '/logs' || to.path === '/cadastros' || to.path === '/solicitacoes') &&
    user.value.level !== 'admin'
  ) {
    return navigateTo('/')
  }
})
