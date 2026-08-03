<script setup lang="ts">
  import type { TableColumn } from '@nuxt/ui'
  import type { Column } from '@tanstack/vue-table'
  import type { VNode } from 'vue'

  import { NuxtTime, UButton } from '#components'

  useHead({ title: 'Logs' })

  const page = ref(1)
  const itemsPerPage = ref(10)
  const filter = ref('')
  const filterDebounced = refDebounced(filter, 300)

  const { data } = await useApi<LogsResponse>('/logs', {
    method: 'GET',
    query: { page, limit: itemsPerPage, filter: filterDebounced },
    default: () => ({ logs: [], total: 0 }),
  })

  watch(filterDebounced, () => (page.value = 1))

  function sortableHeader(column: Column<Log>, label: string): VNode {
    const sorted = column.getIsSorted()
    let icon = 'i-lucide-chevrons-up-down'
    if (sorted === 'asc') icon = 'i-lucide-arrow-up'
    else if (sorted === 'desc') icon = 'i-lucide-arrow-down'

    return h(UButton, {
      color: 'neutral',
      variant: 'ghost',
      size: 'xs',
      label,
      icon,
      class: '-mx-2 font-semibold uppercase tracking-wider text-xs',
      onClick: () => column.toggleSorting(sorted === 'asc'),
    })
  }

  const columns: TableColumn<Log>[] = [
    {
      accessorKey: 'zone',
      header: ({ column }) => sortableHeader(column, 'Registro'),
      cell: ({ row }) =>
        h('span', { class: 'data font-medium text-highlighted' }, row.original.zone),
    },
    {
      accessorKey: 'username',
      header: ({ column }) => sortableHeader(column, 'Usuário'),
      cell: ({ row }) => h('span', { class: 'data text-sm text-toned' }, row.original.username),
    },
    {
      accessorKey: 'details',
      header: ({ column }) => sortableHeader(column, 'Detalhes'),
      cell: ({ row }) => h('span', { class: 'text-sm text-muted' }, row.original.details),
    },
    {
      accessorKey: 'createdAt',
      header: ({ column }) => sortableHeader(column, 'Criado em'),
      cell: ({ row }) =>
        h('span', { class: 'data tnum text-sm whitespace-nowrap text-dimmed' }, [
          h(NuxtTime, {
            datetime: row.original.createdAt,
            day: '2-digit',
            month: '2-digit',
            year: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            locale: 'pt-BR',
          }),
        ]),
    },
  ]
</script>

<template>
  <UContainer class="py-8 sm:py-10">
    <PageHeader
      eyebrow="Auditoria"
      title="Logs do sistema"
      description="Toda alteração feita nas zonas e nos registros, com autor e horário." />

    <div class="overflow-hidden rounded-xl border border-default bg-default">
      <div
        class="flex flex-col gap-3 border-b border-default bg-muted/40 p-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-2.5">
          <UIcon name="i-lucide-scroll-text" class="size-4 shrink-0 text-dimmed" />
          <p class="text-sm font-medium text-toned">
            <span class="data text-highlighted tnum">{{ data.total }}</span>
            {{ data.total === 1 ? 'registro' : 'registros' }}
          </p>
        </div>

        <UInput
          v-model="filter"
          icon="i-lucide-search"
          placeholder="Filtrar por zona, usuário ou detalhe..."
          class="w-full sm:w-80"
          :ui="{ base: 'data' }" />
      </div>

      <UTable
        :data="data.logs"
        :columns
        :ui="{
          tr: 'transition-colors duration-150 hover:bg-elevated/60',
          th: 'py-2',
          td: 'py-2.5',
        }">
        <template #empty>
          <div class="flex flex-col items-center gap-2 py-12 text-center">
            <UIcon name="i-lucide-file-clock" class="size-7 text-dimmed" />
            <p class="text-sm font-medium text-toned">Nenhum log encontrado</p>
            <p class="max-w-xs text-xs text-dimmed">
              {{
                filter
                  ? 'Nenhum resultado para esse filtro. Tente outro termo.'
                  : 'Ainda não há atividade registrada.'
              }}
            </p>
          </div>
        </template>
      </UTable>

      <div
        v-if="data.total > itemsPerPage"
        class="flex justify-center border-t border-default bg-muted/40 p-3">
        <UPagination
          v-model:page="page"
          active-variant="subtle"
          :total="data.total"
          :items-per-page="itemsPerPage" />
      </div>
    </div>
  </UContainer>
</template>
