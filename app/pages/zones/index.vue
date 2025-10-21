<script setup lang="ts">
useHead({ title: 'Dashboard zones' })

useSeoMeta({ description: 'Dashboard zones page' })

defineOgImageComponent('Techs', { title: 'Dashboard zones page' })

const zoneId = ref('')
</script>

<template>
  <UContainer class="p-4">
    <Zones v-if="!zoneId" v-model:zone-id="zoneId" />
    <Record v-else-if="zoneId" v-model:zone-id="zoneId" />
  </UContainer>
</template>
