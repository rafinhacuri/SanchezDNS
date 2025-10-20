<script setup lang="ts">
useHead({ title: 'Dashboard zones' })

useSeoMeta({ description: 'Dashboard zones page' })

defineOgImageComponent('Techs', { title: 'Dashboard zones page' })

const zoneId = ref('')
</script>

<template>
  <UContainer class="p-4">
    <transition-group mode="out-in" enter-active-class="transition-all duration-500 ease-out" enter-from-class="opacity-0 translate-x-10" enter-to-class="opacity-100 translate-x-0">
      <Zones v-if="!zoneId" v-model:zone-id="zoneId" />
      <Record v-else-if="zoneId" v-model:zone-id="zoneId" />
    </transition-group>
  </UContainer>
</template>
