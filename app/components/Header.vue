<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

const route = useRoute()

const { optionSelected, optionsConection } = useConection()

const items = computed<NavigationMenuItem[]>(() => [
  {
    label: 'Zones',
    icon: 'i-lucide-database',
    to: '/zones/dashboard',
    active: route.path === '/zones/dashboard' || route.path.startsWith('/zones/templates'),
    children: [
      {
        label: 'Dashboard',
        icon: 'i-lucide-chart-pie',
        description: 'View and manage all DNS zones.',
        to: '/zones/dashboard',
      },
      {
        label: 'Templates',
        icon: 'i-lucide-file-text',
        description: 'Manage DNS templates.',
        to: '/zones/templates',
      },
    ],
  },
  {
    label: 'Server',
    icon: 'i-lucide-server',
    active: route.path.startsWith('/srv/statistics') || route.path.startsWith('/srv/configuration'),
    children: [
      {
        label: 'Statistics',
        icon: 'i-lucide-bar-chart-3',
        description: 'View server statistics.',
        to: '/srv/statistics',
      },
      {
        label: 'Configuration',
        icon: 'i-lucide-settings',
        description: 'Manage server configuration.',
        to: '/srv/configuration',
      },
    ],
  },
  {
    label: 'Configuration',
    icon: 'i-lucide-settings',
    active: route.path.startsWith('/config/dns-machines') || route.path.startsWith('/config/api-keys'),
    children: [
      {
        label: 'DNS Machines',
        icon: 'i-lucide-computer',
        description: 'Manage DNS machines connections.',
        to: '/config/dns-machines',
      },
      {
        label: 'API Keys',
        icon: 'i-lucide-key',
        description: 'Manage API keys.',
        to: '/config/api-keys',
      },
    ],
  },
  {
    label: 'Users',
    icon: 'i-lucide-users',
    to: '/user/users',
    active: route.path.startsWith('/user/users') || route.path.startsWith('/user/groups'),
    children: [
      {
        label: 'Manage Users',
        icon: 'i-lucide-user',
        description: 'View and manage all users.',
        to: '/user/users',
      },
      {
        label: 'Manage Groups',
        icon: 'i-lucide-users',
        description: 'View and manage all groups.',
        to: '/user/groups',
      },
    ],
  },
  {
    label: 'Logs',
    icon: 'i-lucide-git-fork',
    to: '/logs',
    active: route.path.startsWith('/logs'),
  },
])
</script>

<template>
  <UHeader>
    <template #title>
      <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="32" />
      Sanchez<span class="text-green-500">DNS</span>
    </template>

    <UNavigationMenu v-if="route.path !== '/login'" :items="items" />

    <template #right>
      <USelectMenu v-if="route.path !== '/login'" v-model="optionSelected" :items="optionsConection" placeholder="Select a connection" size="sm" />

      <UTooltip text="Open on GitHub">
        <UButton color="neutral" variant="ghost" to="https://github.com/rafinhacuri/SanchezDNS" target="_blank" icon="i-simple-icons-github" aria-label="GitHub" />
      </UTooltip>
    </template>

    <template #body>
      <UNavigationMenu :items="items" orientation="vertical" class="-mx-2.5" />
    </template>
  </UHeader>
</template>
