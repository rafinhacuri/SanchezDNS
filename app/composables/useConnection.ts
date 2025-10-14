export default createGlobalState(() => {
  const optionsConnection = ref([
    'connection1',
    'connection2',
    'connection3',
  ])
  const optionSelected = ref('connection1')

  watch(optionSelected, async () => {
    await navigateTo('/zones/dashboard')
  })

  return { optionsConnection, optionSelected }
})
