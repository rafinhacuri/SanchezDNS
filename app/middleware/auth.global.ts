export default defineNuxtRouteMiddleware(async to => {
  const { isLoggedIn, setUserSession, clearUserSession } = useUserSession()

  const res = await useRequestFetch()<{ username: string, isAdmin: boolean }>('/server/api/check-session').catch(() => null)

  if(!res){
    clearUserSession()
    if(to.fullPath === '/') return
    return navigateTo('/')
  }
  setUserSession({ username: res.username, admin: res.isAdmin })

  // * Logado tentando acessar login
  if(isLoggedIn.value && to.fullPath === '/') return navigateTo('/zones/dashboard')

  // * Não logado tentando qualquer rota exceto login
  if(!isLoggedIn.value && to.fullPath !== '/') return navigateTo('/')
})
