<script setup lang="ts">
  import type { NavigationMenuItem } from '@nuxt/ui'

  const { user } = useUser()
  const route = useRoute()

  const items = computed<NavigationMenuItem[]>(() => {
    const base: NavigationMenuItem[] = [
      { label: 'Zonas', icon: 'i-lucide-database', to: '/', active: route.path === '/' },
      {
        label: 'Estatísticas',
        icon: 'i-lucide-bar-chart-3',
        to: '/statistics',
        active: route.path === '/statistics',
      },
    ]

    if (user.value.level === 'admin') {
      base.push({
        label: 'Logs',
        icon: 'i-lucide-fingerprint',
        to: '/logs',
        active: route.path === '/logs',
      })
    }

    if (user.value.level === 'admin') {
      base.push({
        label: 'Solicitações',
        icon: 'i-lucide-clipboard-list',
        to: '/solicitacoes',
        active: route.path === '/solicitacoes',
      })
    }

    if (user.value.level === 'admin') {
      base.push({
        label: 'Cadastros',
        icon: 'i-lucide-users',
        to: '/cadastros',
        active: route.path === '/cadastros',
      })
    }

    return base
  })

  const nuxtReady = ref(false)
  onNuxtReady(() => (nuxtReady.value = true))
</script>

<template>
  <UHeader>
    <template #title>
      <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="32" />
      Sanchez<span class="text-green-500">DNS</span>
    </template>

    <UNavigationMenu :items />

    <template #right>
      <div class="flex items-center gap-4">
        <UColorModeSwitch size="lg" />
        <User />
      </div>
    </template>

    <template #body>
      <UNavigationMenu :items orientation="vertical" class="-mx-2.5" />
    </template>
  </UHeader>
</template>
