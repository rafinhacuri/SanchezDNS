<script setup lang="ts">
  useHead({ title: 'Zonas' })

  const route = useRoute()

  const items = [
    { id: 'normal', label: 'Domínios', icon: 'i-lucide-globe' },
    { id: 'reverse', label: 'in-addr.arpa', icon: 'i-lucide-undo-2' },
    { id: 'reverse-ipv6', label: 'ip6.arpa', icon: 'i-lucide-undo-2' },
  ]

  const selectedTab = computed(() => {
    const tipo = String(route.query.type ?? 'normal')
    return items.some((item) => item.id === tipo) ? tipo : 'normal'
  })

  async function selectTab(id: string): Promise<void> {
    await navigateTo({ query: id === 'normal' ? {} : { type: id } }, { replace: true })
  }
</script>

<template>
  <UContainer class="py-8 sm:py-10">
    <PageHeader
      eyebrow="Servidor DNS"
      title="Zonas"
      description="Selecione uma zona para gerenciar seus registros DNS.">
      <template #actions>
        <div class="flex rounded-lg bg-elevated p-1 ring-1 ring-default">
          <button
            v-for="item of items"
            :key="item.id"
            type="button"
            class="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors duration-150"
            :class="
              selectedTab === item.id
                ? 'bg-default text-highlighted shadow-sm ring-1 ring-default'
                : 'text-muted hover:text-highlighted'
            "
            :aria-pressed="selectedTab === item.id"
            @click="selectTab(item.id)">
            <UIcon :name="item.icon" class="size-4 shrink-0" />
            <span :class="item.id === 'normal' ? '' : 'data text-xs'">{{ item.label }}</span>
          </button>
        </div>
      </template>
    </PageHeader>

    <Zones :type="selectedTab" />
  </UContainer>
</template>
