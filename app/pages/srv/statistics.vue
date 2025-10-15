<script setup lang="ts">
useHead({ title: 'Server Statistics' })
useSeoMeta({ description: 'Server statistics overview' })
defineOgImageComponent('Techs', { title: 'Server statistics overview' })

const { optionSelected } = await useConnection()

const { data, error, refresh } = await useFetch<{ zones: number, records: number, uptime: string, users: number, status: string, qps: number, udpQueries: number, tcpQueries: number, serverId: string, startedAt: string }>('/server/api/statistics', { method: 'GET', query: { connection: optionSelected.value } })

if(error.value) throw createError({ statusCode: 500, statusMessage: 'Failed to fetch statistics' })

onNuxtReady(() => setInterval(refresh, 60000))
</script>

<template>
  <UContainer>
    <div v-if="data" class="space-y-8 py-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold">
            Server Statistics
          </h1>
          <p class="text-sm text-gray-500">
            Overview of your DNS server statistics and performance
          </p>
        </div>
        <UBadge variant="soft">
          Online
        </UBadge>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-500">
                Zones
              </div>
              <div class="text-2xl font-semibold">
                {{ data.zones }}
              </div>
            </div>
            <UIcon name="i-lucide-layers" class="h-6 w-6" />
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-500">
                Records
              </div>
              <div class="text-2xl font-semibold">
                {{ data.records }}
              </div>
            </div>
            <UIcon name="i-lucide-database" class="h-6 w-6" />
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-500">
                Users
              </div>
              <div class="text-2xl font-semibold">
                {{ data.users }}
              </div>
            </div>
            <UIcon name="i-lucide-users" class="h-6 w-6" />
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-500">
                Uptime
              </div>
              <div class="text-2xl font-semibold">
                {{ data.uptime }}
              </div>
            </div>
            <UIcon name="i-lucide-timer" class="h-6 w-6" />
          </div>
        </UCard>
      </div>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <div class="font-medium">
              Traffic Statistics
            </div>
          </div>
        </template>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
          <div>
            <div class="text-sm text-gray-500">
              QPS (req/s)
            </div>
            <div class="text-xl font-semibold">
              {{ data.qps }}
            </div>
          </div>
          <div>
            <div class="text-sm text-gray-500">
              UDP Queries
            </div>
            <div class="text-xl font-semibold">
              {{ data.udpQueries.toLocaleString() }}
            </div>
          </div>
          <div>
            <div class="text-sm text-gray-500">
              TCP Queries
            </div>
            <div class="text-xl font-semibold">
              {{ data.tcpQueries.toLocaleString() }}
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="font-medium">
            Server Information
          </div>
        </template>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div>
            <div class="text-sm text-gray-500">
              Server ID
            </div>
            <div class="text-lg font-semibold">
              {{ data.serverId }}
            </div>
          </div>
          <div>
            <div class="text-sm text-gray-500">
              Started At
            </div>
            <div class="text-lg font-semibold">
              {{ new Date(data.startedAt).toLocaleString() }}
            </div>
          </div>
        </div>
      </UCard>
    </div>
    <div v-else class="flex items-center justify-center h-64">
      <UProgress animation="carousel" />
    </div>
  </UContainer>
</template>
