export default createGlobalState(() => {
  const optionsConection = ref([
    'connection1',
    'connection2',
    'connection3',
  ])
  const optionSelected = ref('')

  watch(optionSelected, async () => {
    await navigateTo('/zones/dashboard')
  })

  return { optionsConection, optionSelected }
})
