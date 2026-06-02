<script setup lang="ts">
  useHead({ title: 'Estatísticas' })

  const { data, refresh } = await useFetch<StatisticsResponse>('/server/api/statistics')

  onNuxtReady(() => setInterval(refresh, 60_000))
</script>

<template>
  <UContainer>
    <div v-if="data" class="space-y-8 py-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold">Estatísticas do Servidor</h1>
          <p class="text-sm text-gray-500">
            Visão geral das estatísticas e desempenho servidor DNS
          </p>
        </div>
        <UBadge variant="soft"> Online </UBadge>
      </div>

      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-500">Zonas</div>
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
              <div class="text-sm text-gray-500">Registros</div>
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
              <div class="text-sm text-gray-500">Uptime</div>
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
            <div class="font-medium">Estatísticas de Tráfego</div>
          </div>
        </template>

        <div class="grid grid-cols-1 gap-6 sm:grid-cols-3">
          <div>
            <div class="text-sm text-gray-500">UDP Queries</div>
            <div class="text-xl font-semibold">
              {{ data.udpQueries.toLocaleString() }}
            </div>
          </div>
          <div>
            <div class="text-sm text-gray-500">TCP Queries</div>
            <div class="text-xl font-semibold">
              {{ data.tcpQueries.toLocaleString() }}
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="font-medium">Informações do Servidor</div>
        </template>

        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <div>
            <div class="text-sm text-gray-500">ID do Servidor</div>
            <div class="text-lg font-semibold">ns1.cbpf.br</div>
          </div>
          <div>
            <div class="text-sm text-gray-500">Iniciado Em</div>
            <NuxtTime
              class="text-lg font-semibold"
              :datetime="data.startedAt"
              month="2-digit"
              day="2-digit"
              year="numeric"
              hour="2-digit"
              minute="2-digit"
              second="2-digit" />
          </div>
        </div>
      </UCard>
    </div>
  </UContainer>
</template>
