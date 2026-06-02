<script setup lang="ts">
  import type { TableColumn } from '@nuxt/ui'

  import { NuxtTime, UButton } from '#components'

  useHead({ title: 'Logs' })

  const page = ref(1)
  const itemsPerPage = ref(10)
  const filter = ref('')
  const filterDebounced = refDebounced(filter, 300)

  const { data } = await useFetch<LogsResponse>('/server/api/logs', {
    method: 'GET',
    query: { page, limit: itemsPerPage, filter: filterDebounced },
    default: () => ({ logs: [], total: 0 }),
  })

  watch(filterDebounced, () => (page.value = 1))

  const columns: TableColumn<Log>[] = [
    {
      accessorKey: 'zone',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Zona',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
    },
    {
      accessorKey: 'username',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Usuário',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
    },
    {
      accessorKey: 'details',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Detalhes',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
    },
    {
      accessorKey: 'createdAt',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Criado em',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
      cell: ({ row }) =>
        h(NuxtTime, {
          datetime: row.original.createdAt,
          day: '2-digit',
          month: '2-digit',
          year: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
          locale: 'pt-BR',
        }),
    },
  ]
</script>

<template>
  <UContainer class="py-8">
    <div class="my-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">Logs do Sistema</h1>
        <p class="text-gray-500 dark:text-gray-400">
          Monitore todas as atividades do sistema e ações dos usuários.
        </p>
      </div>
      <UInput v-model="filter" placeholder="Pesquisar logs..." class="mb-4" />
    </div>

    <UTable :data="data.logs" :columns />

    <div v-if="data.total > itemsPerPage" class="flex justify-center border-t border-default pt-4">
      <UPagination
        v-model:page="page"
        active-color="info"
        color="info"
        active-variant="subtle"
        :total="data.total"
        :items-per-page="itemsPerPage" />
    </div>
  </UContainer>
</template>
