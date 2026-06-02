async function useLogout(): Promise<GoRes> {
  const res = await $fetch<GoRes>('/server/api/logout', {
    method: 'post',
  })
  return res
}

export function useUser(): {
  user: ComputedRef<{ idcbpf: string; representar: boolean; level: string }>
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
    idcbpf: data.value?.idcbpf ?? '',
    representar: data.value?.representar ?? false,
    level: data.value?.level ?? '',
  }))

  async function refresh(): Promise<void> {
    await refreshU()
  }

  const isLoggedIn = computed(() => Boolean(user.value.idcbpf))

  return { user, isLoggedIn, refresh, useLogout }
}
