<script setup lang="ts">
  import type { DropdownMenuItem } from '@nuxt/ui'

  const { user, useLogout } = useUser()

  const itemsDropdown = computed<DropdownMenuItem[][]>(() => [
    [
      {
        label: user.value.email,
        icon: 'i-lucide-user',
        onClick: async (): Promise<void> => {
          await navigateTo('/perfil')
        },
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
  <UDropdownMenu
    :items="itemsDropdown"
    :ui="{
      content:
        'w-64 rounded-2xl border border-default bg-default/95 backdrop-blur shadow-lg ring-1 ring-default/60',
      item: 'flex items-center gap-2 rounded-xl px-3 py-2 text-sm font-medium transition hover:bg-info/10',
      label: 'text-xs font-semibold text-muted uppercase px-3 pt-2 pb-1',
      separator: 'my-2 border-default',
    }">
    <button
      type="button"
      class="group relative grid size-10 place-items-center rounded-full ring-1 ring-gray-200/70 transition hover:scale-[1.02] hover:ring-blue-500/30 dark:ring-gray-800/70"
      :aria-label="`Menu do usuário ${user.email}`">
      <img
        :src="`/server/api/file/${user.email}`"
        :alt="user.email"
        class="size-9 rounded-full object-cover shadow-sm" />
      <span
        class="pointer-events-none absolute -right-0.5 -bottom-0.5 size-3 rounded-full bg-blue-300 ring-2 ring-white dark:ring-gray-950" />
    </button>
  </UDropdownMenu>
</template>
