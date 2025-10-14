<script setup lang="ts">
useHead({ title: 'Dashboard zones' })

useSeoMeta({ description: 'Dashboard zones page' })

defineOgImageComponent('Techs', { title: 'Dashboard zones page' })
const { user } = useUserSession()
</script>

<template>
  <p>{{ user }}</p>
</template>
