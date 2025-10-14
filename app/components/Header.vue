<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

const route = useRoute()

const items = computed<NavigationMenuItem[]>(() => [
  {
    label: 'Inicio',
    to: '/',
    active: route.path.startsWith('/'),
  },
])
</script>

<template>
  <UHeader>
    <template #title>
      Sanchez<span class="text-green-500">DNS</span>
    </template>

    <UNavigationMenu :items="items" />

    <template #right>
      <UTooltip text="Open on GitHub">
        <UButton color="neutral" variant="ghost" to="https://github.com/rafinhacuri/SanchezDNS" target="_blank" icon="i-simple-icons-github" aria-label="GitHub" />
      </UTooltip>
    </template>
  </UHeader>
</template>
