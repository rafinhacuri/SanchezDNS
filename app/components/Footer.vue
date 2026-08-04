<script setup lang="ts">
  const { version } = useRuntimeConfig().public

  const currentYear = new Date().getFullYear()
</script>

<template>
  <UFooter :ui="{ root: 'border-t border-default bg-muted/40' }">
    <template #left>
      <div class="space-y-1">
        <p class="text-xs text-muted">
          ©&nbsp;{{ currentYear }}
          <span class="font-medium text-toned">SanchezDNS</span>
        </p>
        <p class="text-[11px] text-dimmed">
          Gerenciador de servidores DNS autoritativos com PowerDNS
        </p>
      </div>
    </template>

    <template #right>
      <span
        class="rounded-full bg-elevated px-2 py-1 data text-[10px] leading-none font-medium text-dimmed ring-1 ring-default">
        v{{ version }}
      </span>
    </template>
  </UFooter>
</template>
