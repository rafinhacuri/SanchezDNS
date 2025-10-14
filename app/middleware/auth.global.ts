export default defineNuxtRouteMiddleware(async to => {
  const { isLoggedIn } = useUserSession()

  // * Logado tentando acessar login
  if(isLoggedIn.value && ['/login'].includes(to.fullPath)) return navigateTo('/zones/dashboard')

  // * Não logado tentando qualquer rota exceto login
  if(!isLoggedIn.value && !['/login'].includes(to.fullPath)) return navigateTo('/login')
})
