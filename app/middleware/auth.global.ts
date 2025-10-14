export default defineNuxtRouteMiddleware(async to => {
  const { isLoggedIn, setUserSession, clearUserSession, user } = useUserSession()
  const { optionSelected } = await useConnection()

  const res = await useRequestFetch()<{ username: string, isAdmin: boolean }>('/server/api/check-session').catch(() => null)

  if(!res){
    clearUserSession()
    if(to.fullPath !== '/') return navigateTo('/')
    return
  }
  setUserSession({ username: res.username, admin: res.isAdmin })

  // * Logado tentando acessar login
  if(isLoggedIn.value && to.fullPath === '/') return navigateTo('/zones/dashboard')

  // * Não logado tentando qualquer rota exceto login
  if(!isLoggedIn.value && to.fullPath !== '/') return navigateTo('/')

  // * logado tentando acessar rota admin
  if(!user.value?.admin && to.fullPath.startsWith('/users')) return navigateTo('/zones/dashboard')

  // * nao selecionou conexao
  if(!optionSelected.value && !['/', '/start'].includes(to.fullPath) && isLoggedIn) return navigateTo('/start')
  if(optionSelected.value && ['/start'].includes(to.fullPath) && isLoggedIn) return navigateTo('/zones/dashboard')
})
