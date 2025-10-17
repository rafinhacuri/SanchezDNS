<script setup lang="ts">
useHead({ title: 'Dashboard zones' })

useSeoMeta({ description: 'Dashboard zones page' })

defineOgImageComponent('Techs', { title: 'Dashboard zones page' })

const zoneId = ref('')
</script>

<template>
  <UContainer class="p-4">
    <Zones v-model:zone-id="zoneId" />
  </UContainer>
</template>
