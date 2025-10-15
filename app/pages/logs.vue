<script setup lang="ts">
useHead({ title: 'Logs' })
useSeoMeta({ description: 'Logs page' })
defineOgImageComponent('Techs', { title: 'Logs page' })

interface Log{
  id: string,
  idConnection: string,
  hostServer: string,
  zone: string,
  username: string,
  ip: string,
  action: string,
  details: string,
  createdAt: number,
}

const { data } = await useFetch<Log[]>('/server/api/logs', { method: 'GET' })
</script>

<template>
  <UContainer>
    <div class="flex items-center justify-between my-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">
          System Logs
        </h1>
        <p class="text-gray-500 dark:text-gray-400">
          Monitor all system activities and user actions.
        </p>
      </div>
    </div>
    {{ data }}
  </UContainer>
</template>
