<script setup lang="ts">
  import { pt_br } from '@nuxt/ui/locale'
</script>

<template>
  <UApp :toaster="{ expand: false }" :locale="pt_br">
    <NuxtRouteAnnouncer />
    <NuxtLoadingIndicator />
    <NuxtLayout>
      <NuxtPage :transition="{ name: 'page', mode: 'out-in' }" />
    </NuxtLayout>
  </UApp>
</template>
