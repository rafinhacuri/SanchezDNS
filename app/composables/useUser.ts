// oxlint-disable-next-line harlanzw/vue-no-faux-composables
async function useLogout(): Promise<GoRes> {
  const res = await $fetch<GoRes>('/server/api/logout', {
    method: 'post',
  })
  return res
}

export function useUser(): {
  user: ComputedRef<{ email: string; level: string }>
  isLoggedIn: ComputedRef<boolean>
  refresh: () => Promise<void>
  useLogout: () => Promise<GoRes>
} {
  const headers = useRequestHeaders(['cookie'])

  const { data, refresh: refreshU } = useFetch<SessionRes>('/server/api/session', {
    immediate: false,
    credentials: 'include',
    headers,
  })

  const user = computed(() => ({
    email: data.value?.email ?? '',
    level: data.value?.level ?? '',
  }))

  async function refresh(): Promise<void> {
    await refreshU()
  }

  const isLoggedIn = computed(() => Boolean(user.value.email))

  return { user, isLoggedIn, refresh, useLogout }
}
