<template>
  <UApp :toaster="{expand: false}">
    <NuxtRouteAnnouncer />
    <NuxtLoadingIndicator color="repeating-linear-gradient(to right,#1d4ed8 0%,#60a5fa 100%)" />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>
