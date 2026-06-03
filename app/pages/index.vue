<script setup lang="ts">
  import type { TabsItem } from '@nuxt/ui'

  useHead({ title: 'Zonas' })

  const selectedTab = ref('normal')

  const items = ref<TabsItem[]>([
    {
      id: 'normal',
      label: 'Zonas',
      icon: 'i-lucide-database',
    },
    {
      id: 'reverse',
      label: 'Zonas in-addr.arpa',
      icon: 'i-lucide-database',
    },
    {
      id: 'reverse-ipv6',
      label: 'Zonas ip6.arpa',
      icon: 'i-lucide-database',
    },
  ])

  const zoneId = ref('')
  const nivel = ref('')
</script>

<template>
  <UContainer class="py-8">
    <template v-if="!zoneId">
      <UTabs v-model="selectedTab" :items value-key="id" />
      <div class="mt-6">
        <Zones v-model:zone-id="zoneId" v-model:type="selectedTab" v-model:nivel="nivel" />
      </div>
    </template>

    <Record v-else v-model:zone-id="zoneId" v-model:nivel="nivel" />
  </UContainer>
</template>
