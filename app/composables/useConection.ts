export default createGlobalState(() => {
  const optionsConection = ref([
    'connection1',
    'connection2',
    'connection3',
  ])
  const optionSelected = ref(optionsConection.value[0])

  return {
    optionsConection,
    optionSelected,
  }
})
