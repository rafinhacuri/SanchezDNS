// oxlint-disable typescript/explicit-function-return-type typescript/explicit-module-boundary-types no-underscore-dangle typescript/no-unsafe-assignment typescript/no-unsafe-member-access harlanzw/vue-require-composable-prefix typescript/no-unsafe-return
export function $api<T = GoRes>(
  url: string,
  options?: Parameters<ReturnType<typeof useNuxtApp>['$api']>[1],
) {
  const nuxtApp = useNuxtApp()

  return nuxtApp.$api<T>(url, {
    onResponse({ response }) {
      const data = response._data

      if (typeof response._data === 'string') {
        response._data = { message: data }
      }
    },
    ...options,
  })
}

export const useApi = createUseFetch((callerOptions) => {
  const headers = useRequestHeaders([
    'cookie',
    'authorization',
    'x-forwarded-for',
    'x-real-ip',
    'x-forwarded-proto',
    'x-forwarded-host',
    'x-forwarded-port',
    'user-agent',
    'accept-language',
    'referer',
    'x-request-id',
    'traceparent',
    'tracestate',
  ])

  return {
    $fetch: useNuxtApp().$api,
    headers,
    credentials: 'include',
    ...callerOptions,
  }
})
