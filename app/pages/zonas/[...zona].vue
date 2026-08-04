<script setup lang="ts">
  const route = useRoute()

  const zona = computed(() => {
    const param = route.params.zona
    return Array.isArray(param) ? param.join('/') : String(param ?? '')
  })

  useHead({ title: computed(() => zona.value) })
</script>

<template>
  <UContainer class="py-8 sm:py-10">
    <Record :zone="zona" />
  </UContainer>
</template>
