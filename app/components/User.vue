<script setup lang="ts">
  import type { DropdownMenuItem } from '@nuxt/ui'

  const { user, useLogout } = useUser()

  const baseUrl = useApiUrl()

  const itemsDropdown = computed<DropdownMenuItem[][]>(() => [
    [
      {
        label: user.value.email,
        icon: 'i-lucide-user',
      },
      {
        label: user.value.level === 'admin' ? 'Administrador' : 'Usuário',
        icon: 'i-lucide-shield',
      },
    ],
    [
      {
        label: 'Logout',
        icon: 'i-lucide-log-out',
        onClick: useLogout,
      },
    ],
  ])
</script>

<template>
  <UDropdownMenu :items="itemsDropdown" :content="{ align: 'end' }">
    <UButton
      variant="ghost"
      color="neutral"
      class="gap-2 px-1.5"
      :aria-label="`Conta de ${user.email}`">
      <img
        :src="`${baseUrl}/file/${user.email}`"
        :alt="user.email"
        class="size-7 rounded-full object-cover ring-1 ring-default" />
      <span class="hidden max-w-40 truncate text-sm font-medium lg:inline">{{ user.email }}</span>
      <UIcon name="i-lucide-chevron-down" class="hidden size-4 text-dimmed lg:inline" />
    </UButton>
  </UDropdownMenu>
</template>
