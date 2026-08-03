export default defineNuxtPlugin(() => {
  const baseURL = useApiUrl()

  const api = $fetch.create({
    baseURL,
    onResponseError({ response }) {
      if (response.status === 403) {
        globalThis.location.reload()
      }
    },
  })

  return {
    provide: {
      api,
    },
  }
})
