export default defineNuxtPlugin(() => {
  const { siteUrl } = useRuntimeConfig().public

  const baseURL = `${siteUrl.replace(/\/$/u, '')}/go`

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
