export function useApiUrl(): string {
  const { siteUrl } = useRuntimeConfig().public

  return `${siteUrl.replace(/\/$/u, '')}/go`
}
