<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'

useHead({ title: 'Logs' })
useSeoMeta({ description: 'Logs page' })
defineOgImageComponent('Techs', { title: 'Logs page' })

interface Log{
  id: string,
  idConnection: string,
  hostServer: string,
  zone: string,
  username: string,
  action: string,
  details: string,
  createdAt: Date,
}

const page = ref(1)
const itemsPerPage = ref(10)

const globalFilter = ref('')

const { data } = await useFetch<{ data: Log[], total: number }>('/server/api/logs', { method: 'GET', query: { page, limit: itemsPerPage, filter: globalFilter } })

const table = useTemplateRef('table')

watch(globalFilter, () => {
  page.value = 1
})

const UButton = resolveComponent('UButton')
const NuxtTime = resolveComponent('NuxtTime')

const columns: TableColumn<Log>[] = [
  {
    accessorKey: 'hostServer',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Host Server', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'zone',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Zone', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'username',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Username', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'details',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Details', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Created At', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
    cell: ({ row }) => h(NuxtTime, { datatime: row.original.createdAt, day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' }),
  },
]
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
      <UInput v-model="globalFilter" placeholder="Search logs..." class="mb-4" />
    </div>

    <UTable ref="table" :data="data?.data" :columns="columns" />

    <div v-if="data" class="border-default flex justify-center border-t pt-4">
      <UPagination v-if="data.total > itemsPerPage" v-model:page="page" active-color="primary" active-variant="subtle" :total="data.total" :items-per-page="itemsPerPage" />
    </div>
  </UContainer>
</template>
