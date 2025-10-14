<script setup lang="ts">
const { user } = useUserSession()
</script>

<template>
  <div class="flex flex-col items-center justify-center h-full text-center p-4 space-y-4">
    <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="92" class="mt-16" />
    <h1 class="text-2xl font-bold">
      Welcome to SanchezDNS
    </h1>
    <p class="text-lg">
      Please select a DNS machine connection to get started <strong v-if="user.admin">or create a new one</strong>.
    </p>
    <UButton v-if="user.admin" color="primary" variant="outline">
      Create Connection
    </UButton>
  </div>
</template>
