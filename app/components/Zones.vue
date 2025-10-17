<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import type { Row } from '@tanstack/vue-table'
import { getPaginationRowModel } from '@tanstack/vue-table'

const { optionSelected, nameServer } = await useConnection()
const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()

interface Zones{
  name: string,
  id: string,
  serial: number,
}

const { data, refresh } = await useFetch<{ zones: Zones[], message: string }>('/server/api/zones', { method: 'GET', query: { connection: optionSelected } })

if(data.value?.message) throw createError({ statusCode: 500, statusMessage: data.value.message })

const globalFilter = ref('')

const table = useTemplateRef('table')

const pagination = ref({ pageIndex: 0, pageSize: 10 })

watch(globalFilter, () => {
  pagination.value.pageIndex = 0
})

const UButton = resolveComponent('UButton')
const UDropdownMenu = resolveComponent('UDropdownMenu')

const columns: TableColumn<Zones>[] = [
  {
    accessorKey: 'name',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Name', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'serial',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Serial', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    id: 'actions',
    cell: ({ row }) => h('div', { class: 'text-right' }, h(UDropdownMenu, { 'content': { align: 'end' }, 'items': getRowItems(row), 'aria-label': 'Actions dropdown' }, () => h(UButton, { 'icon': 'i-lucide-ellipsis-vertical', 'color': 'neutral', 'variant': 'ghost', 'class': 'ml-auto', 'aria-label': 'Actions dropdown' }))),
  },
]

const { copy } = useClipboard()

function getRowItems(row: Row<Zones>){
  return [
    { type: 'label', label: `${row.original.id} actions` },
    {
      label: 'Copy zone ID',
      icon: 'i-lucide-copy',
      onSelect(){
        copy(row.original.id)
        toast.add({ title: 'Zone ID copied to clipboard!', color: 'success', icon: 'i-lucide-circle-check' })
      },
    },
    { type: 'separator' },
    { label: 'View zone', icon: 'i-lucide-eye' },
    { label: 'Edit zone', icon: 'i-lucide-pencil' },
    { label: 'Delete zone', icon: 'i-lucide-trash', color: 'error' },
  ]
}

const modal = ref(false)

const state = ref<ZoneSchemaType>({ domain: '', soa: { startOfAuthority: '', email: '', refresh: 3600, retry: 600, expire: 604800, negativeCacheTtl: 86400 } })

watch(modal, nv => {
  if(!nv){
    state.value = { domain: '', soa: { startOfAuthority: '', email: '', refresh: 3600, retry: 600, expire: 604800, negativeCacheTtl: 86400 } }
  }
})

async function createZone(){
  start()

  const body = ZoneSchema.safeParse(state.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/zone', { method: 'PUT', body: body.data, query: { connection: optionSelected.value } })
    .catch(error => {
      console.error(error)
      toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
    })

  if(!res) return finish({ error: true })

  if(!res.message){
    toast.add({ title: 'An unknown error occurred', icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  finish({ force: true })
  refresh()
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  modal.value = false
}
</script>

<template>
  <div class="flex items-center justify-between my-6">
    <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">
      Zones
    </h1>
    <div class="flex items-center gap-4">
      <UButton label="Create Zone" icon="i-lucide-plus" variant="soft" @click="modal = true" />
      <UInput v-model="globalFilter" placeholder="Search zones..." />
    </div>
  </div>

  <UTable ref="table" v-model:global-filter="globalFilter" v-model:pagination="pagination" :pagination-options="{ getPaginationRowModel: getPaginationRowModel()}" :data="data?.zones" :columns="columns" />

  <div class="border-default flex justify-center border-t pt-4">
    <UPagination :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1" :items-per-page="table?.tableApi?.getState().pagination.pageSize" :total="table?.tableApi?.getFilteredRowModel().rows.length" @update:page="(p) => table?.tableApi?.setPageIndex(p - 1)" />
  </div>

  <UModal v-model:open="modal" title="Create Zone" :description="`Create a new zone in ${nameServer}`" :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :schema="ZoneSchema" :state="state" class="space-y-4">
        <UFormField label="Domain" name="domain">
          <UInput v-model="state.domain" icon="i-lucide-computer" class="w-full" placeholder="Ex: example.com" />
        </UFormField>

        <USeparator label="Start of Authority (SOA) Record Settings" />

        <UFormField label="Start of Authority" name="soa.startOfAuthority">
          <UInput v-model="state.soa.startOfAuthority" icon="i-lucide-shield-check" class="w-full" placeholder="Ex: ns1.example.com" />
        </UFormField>
        <UFormField label="Email" name="soa.email">
          <UInput v-model="state.soa.email" icon="i-lucide-mail" class="w-full" placeholder="Ex: hostmaster.example.com" />
        </UFormField>
        <UFormField label="Refresh" name="soa.refresh">
          <UInput v-model="state.soa.refresh" type="number" icon="i-lucide-refresh-cw" class="w-full" placeholder="Ex: 3600" />
        </UFormField>
        <UFormField label="Retry" name="soa.retry">
          <UInput v-model="state.soa.retry" type="number" icon="i-lucide-clock" class="w-full" placeholder="Ex: 600" />
        </UFormField>
        <UFormField label="Expire" name="soa.expire">
          <UInput v-model="state.soa.expire" type="number" icon="i-lucide-hourglass" class="w-full" placeholder="Ex: 604800" />
        </UFormField>
        <UFormField label="Negative Cache TTL" name="soa.negativeCacheTtl">
          <UInput v-model="state.soa.negativeCacheTtl" type="number" icon="i-lucide-timer" class="w-full" placeholder="Ex: 3600" />
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
      <UButton label="Confirm" :loading="isLoading" @click="createZone" />
    </template>
  </UModal>
</template>
