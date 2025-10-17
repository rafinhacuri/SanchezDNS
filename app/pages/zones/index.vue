<script setup lang="ts">
useHead({ title: 'Dashboard zones' })

useSeoMeta({ description: 'Dashboard zones page' })

defineOgImageComponent('Techs', { title: 'Dashboard zones page' })
</script>

<template>
  <UContainer class="p-4">
    <Zones />
  </UContainer>
</template>
