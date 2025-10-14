<script setup lang="ts">
const colorMode = useColorMode()

onMounted(() => {
  colorMode.value = 'dark'
})
</script>

<template>
  <UApp :toaster="{expand: false}">
    <NuxtRouteAnnouncer />
    <NuxtLoadingIndicator color="repeating-linear-gradient(to right,#1d4ed8 0%,#60a5fa 100%)" />
    <NuxtLayout>
      <NuxtPage :transition="{ name: 'page', mode: 'out-in' }" />
    </NuxtLayout>
  </UApp>
</template>

<style scoped>
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
    filter: blur(2px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
    filter: blur(0);
  }
}

@keyframes fadeOut {
  from {
    opacity: 1;
    transform: translateY(0);
    filter: blur(0);
  }
  to {
    opacity: 0;
    transform: translateY(10px);
    filter: blur(2px);
  }
}

.page-enter-active {
  animation: fadeIn 0.4s ease-out;
}

.page-leave-active {
  animation: fadeOut 0.4s ease-in;
}
</style>
