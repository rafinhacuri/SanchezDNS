<script setup lang="ts">
  import type { NavigationMenuItem } from '@nuxt/ui'

  const { version } = useRuntimeConfig().public

  const { user } = useUser()
  const route = useRoute()

  const items = computed<NavigationMenuItem[]>(() => {
    const base: NavigationMenuItem[] = [
      {
        label: 'Zonas',
        icon: 'i-lucide-database',
        to: '/',
        active: route.path === '/' || route.path.startsWith('/zonas'),
      },
      {
        label: 'Estatísticas',
        icon: 'i-lucide-activity',
        to: '/statistics',
        active: route.path === '/statistics',
      },
    ]

    if (user.value.level === 'admin') {
      base.push(
        {
          label: 'Logs',
          icon: 'i-lucide-scroll-text',
          to: '/logs',
          active: route.path === '/logs',
        },
        {
          label: 'Solicitações',
          icon: 'i-lucide-clipboard-list',
          to: '/solicitacoes',
          active: route.path === '/solicitacoes',
        },
        {
          label: 'Cadastros',
          icon: 'i-lucide-users',
          to: '/cadastros',
          active: route.path === '/cadastros',
        },
      )
    }

    return base
  })
</script>

<template>
  <UHeader :ui="{ root: 'border-b border-default bg-default/80 backdrop-blur-xl' }">
    <template #title>
      <div class="flex items-center gap-2.5">
        <NuxtImg src="/logo.png" alt="SanchezDNS" width="30" height="30" class="shrink-0" />

        <span class="text-base leading-none font-bold tracking-tight text-highlighted">
          Sanchez<span class="text-primary">DNS</span>
        </span>

        <span
          class="hidden rounded-full bg-elevated px-1.5 py-0.5 data text-[10px] leading-none font-medium text-dimmed ring-1 ring-default sm:inline-block">
          v{{ version }}
        </span>
      </div>
    </template>

    <UNavigationMenu :items :ui="{ link: 'font-medium' }" />

    <template #right>
      <div class="flex items-center gap-1.5">
        <UColorModeSwitch />

        <div class="mx-1 hidden h-5 w-px bg-accented sm:block" />

        <User />
      </div>
    </template>

    <template #body>
      <UNavigationMenu :items orientation="vertical" class="-mx-2.5" />
    </template>
  </UHeader>
</template>
