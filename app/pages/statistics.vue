<script setup lang="ts">
  useHead({ title: 'Estatísticas' })

  const { data, refresh } = await useApi<StatisticsResponse>('/statistics')

  onNuxtReady(() => setInterval(refresh, 60_000))

  const kpis = computed(() => [
    {
      label: 'Zonas',
      value: data.value?.zones?.toLocaleString('pt-BR') ?? '0',
      icon: 'i-lucide-globe',
      tint: 'text-blue-500',
    },
    {
      label: 'Registros',
      value: data.value?.records?.toLocaleString('pt-BR') ?? '0',
      icon: 'i-lucide-database',
      tint: 'text-violet-500',
    },
    {
      label: 'Uptime',
      value: data.value?.uptime ?? '—',
      icon: 'i-lucide-timer',
      tint: 'text-emerald-500',
    },
  ])

  const queries = computed(() => [
    {
      label: 'UDP Queries',
      value: data.value?.udpQueries?.toLocaleString('pt-BR') ?? '0',
      icon: 'i-lucide-zap',
    },
    {
      label: 'TCP Queries',
      value: data.value?.tcpQueries?.toLocaleString('pt-BR') ?? '0',
      icon: 'i-lucide-cable',
    },
  ])
</script>

<template>
  <UContainer class="py-8 sm:py-10">
    <PageHeader
      eyebrow="Servidor DNS"
      title="Estatísticas"
      description="Visão geral do desempenho e da carga do servidor. Atualiza a cada minuto.">
      <template #actions>
        <div
          class="inline-flex items-center gap-2 rounded-full bg-emerald-500/10 px-3 py-1.5 ring-1 ring-emerald-500/25">
          <span class="relative flex size-2">
            <span
              class="absolute inline-flex size-full animate-ping rounded-full bg-emerald-500 opacity-60" />
            <span class="relative inline-flex size-2 rounded-full bg-emerald-500" />
          </span>
          <span class="text-xs font-semibold text-emerald-600 dark:text-emerald-400">Online</span>
        </div>
      </template>
    </PageHeader>

    <div v-if="data" class="space-y-4">
      <div class="stagger grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="kpi of kpis"
          :key="kpi.label"
          class="console-rail overflow-hidden rounded-xl border border-default bg-default p-5">
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 space-y-1">
              <p class="text-xs font-semibold tracking-wider text-dimmed uppercase">
                {{ kpi.label }}
              </p>
              <p class="truncate text-3xl font-bold text-highlighted tnum">{{ kpi.value }}</p>
            </div>
            <UIcon :name="kpi.icon" class="size-5 shrink-0" :class="kpi.tint" />
          </div>
        </div>
      </div>

      <div class="grid gap-4 lg:grid-cols-2">
        <div class="overflow-hidden rounded-xl border border-default bg-default">
          <div class="border-b border-default bg-muted/40 px-5 py-3">
            <h2 class="text-sm font-semibold text-highlighted">Tráfego</h2>
          </div>

          <dl class="divide-y divide-default">
            <div
              v-for="q of queries"
              :key="q.label"
              class="flex items-center justify-between gap-4 px-5 py-4">
              <dt class="flex items-center gap-2.5 text-sm text-muted">
                <UIcon :name="q.icon" class="size-4 shrink-0 text-dimmed" />
                {{ q.label }}
              </dt>
              <dd class="data text-lg font-semibold text-highlighted tnum">{{ q.value }}</dd>
            </div>
          </dl>
        </div>

        <div class="overflow-hidden rounded-xl border border-default bg-default">
          <div class="border-b border-default bg-muted/40 px-5 py-3">
            <h2 class="text-sm font-semibold text-highlighted">Servidor</h2>
          </div>

          <dl class="divide-y divide-default">
            <div class="flex items-center justify-between gap-4 px-5 py-4">
              <dt class="flex items-center gap-2.5 text-sm text-muted">
                <UIcon name="i-lucide-server" class="size-4 shrink-0 text-dimmed" />
                ID do servidor
              </dt>
              <dd class="data text-sm font-semibold text-highlighted">{{ data.serverId }}</dd>
            </div>

            <div class="flex items-center justify-between gap-4 px-5 py-4">
              <dt class="flex items-center gap-2.5 text-sm text-muted">
                <UIcon name="i-lucide-power" class="size-4 shrink-0 text-dimmed" />
                Iniciado em
              </dt>
              <dd class="data text-sm font-semibold text-highlighted tnum">
                <NuxtTime
                  :datetime="data.startedAt"
                  month="2-digit"
                  day="2-digit"
                  year="numeric"
                  hour="2-digit"
                  minute="2-digit"
                  second="2-digit"
                  locale="pt-BR" />
              </dd>
            </div>
          </dl>
        </div>
      </div>
    </div>
  </UContainer>
</template>
