export default createGlobalState(async () => {
  const { data: optionsConnection, refresh: refreshConnections } = await useFetch<{ _id: string, name: string }[]>('/server/api/connections', { method: 'GET' })

  const optionSelected = ref('')

  const nameServer = computed(() => optionsConnection.value?.find(option => option._id === optionSelected.value)?.name || '')

  watch(optionSelected, async () => {
    await navigateTo('/zones/dashboard')
  })

  return { optionsConnection, optionSelected, refreshConnections, nameServer }
})
