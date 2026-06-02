<script setup lang="ts">
  import type { DropdownMenuItem, NavigationMenuItem } from '@nuxt/ui'

  const { user, logout } = useUser()
  const route = useRoute()
  const colorMode = useColorMode()

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

    return base
  })

  const nuxtReady = ref(false)
  onNuxtReady(() => (nuxtReady.value = true))
  const queryTheme = computed(() => (nuxtReady.value ? colorMode.value : ''))
</script>

<template>
  <UHeader>
   <template #title>
      <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="32" />
      Sanchez<span class="text-green-500">DNS</span>
    </template>

    <UNavigationMenu :items color="info" />

    <template #right>
      <div class="flex items-center gap-4">
        <UColorModeSwitch size="lg" color="info" />
      </div>
    </template>

    <template #body>
      <UNavigationMenu :items orientation="vertical" class="-mx-2.5" color="info" />
    </template>
  </UHeader>
</template>
