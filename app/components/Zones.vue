<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import type { Row } from '@tanstack/vue-table'
import { getPaginationRowModel } from '@tanstack/vue-table'

const { optionSelected, nameServer } = await useConnection()
const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()

interface Zones{
  name: string,
  kind: string,
  id: string,
  serial: number,
  url: string,
  soa_edit_api: string,
}

const { data } = await useFetch<{ zones: Zones[], message: string }>('/server/api/zones', { method: 'GET', query: { connection: optionSelected } })

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
    accessorKey: 'kind',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Kind', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
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
    accessorKey: 'url',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'URL', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'soa_edit_api',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'SOA Edit API', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
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
      onSelect(){
        copy(row.original.id)
        toast.add({ title: 'Zone ID copied to clipboard!', color: 'success', icon: 'i-lucide-circle-check' })
      },
    },
    { type: 'separator' },
    { label: 'Edit zone' },
    { label: 'View zone' },
  ]
}

const modal = ref(false)

const types = ['Native', 'Primary', 'Secondary']
const soaEditApis = ['DEFAULT', 'INCREASE', 'EPOCH', 'OFF']

const state = ref<ZoneSchemaType>({ domain: '', kind: 'Native', soa_edit_api: 'DEFAULT', masters: [], also_notify: [] })

const typeDetails = computed(() => {
  switch (state.value.kind){
    case 'Native':
      return 'Native - SanchezDNS will not perform any special handling for this zone type Primary or Secondary zone functions.'
    case 'Primary':
      return 'Primary - SanchezDNS will serve as the Primary and will send zone transfers (AXFRs) to other servers configured as Secondaries.'
    case 'Secondary':
      return 'Secondary - SanchezDNS will serve as the Secondary and will request and receive zone transfers (AXFRs) from other servers configured as Primaries.'
    default:
      return ''
  }
})

const soaDetails = computed(() => {
  switch (state.value.soa_edit_api){
    case 'DEFAULT':
      return 'DEFAULT - Generate a SOA serial of YYYYMMDD01. If the current serial is lower than the generated serial, use the generated serial. If the current serial is higher or equal to the generated serial, increase the current serial by 1.'
    case 'INCREASE':
      return 'INCREASE - Increase the current serial by 1.'
    case 'EPOCH':
      return 'EPOCH - Change the serial to the number of seconds since the EPOCH (Unix time).'
    case 'OFF':
      return 'OFF - Disable automatic updates of the SOA serial.'
    default:
      return ''
  }
})

watch(modal, nv => {
  if(!nv){
    state.value = { domain: '', kind: 'Native', soa_edit_api: 'DEFAULT', masters: [], also_notify: [] }
  }
})

async function createZone(){
  start()

  const body = ZoneSchema.safeParse(state.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/zone', { method: 'PUT', body: body.data, query: { connection: optionSelected } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
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
        <UFormField label="Kind" name="kind">
          <div class="flex items-center gap-2">
            <USelect v-model="state.kind" icon="i-lucide-tag" :items="types" class="w-full" />
            <UPopover arrow>
              <UButton icon="i-lucide-info" color="neutral" variant="subtle" class="rounded-full" />

              <template #content>
                <p class="max-w-xs text-sm break-words">
                  {{ typeDetails }}
                </p>
              </template>
            </UPopover>
          </div>
        </UFormField>
        <UFormField v-if="state.kind === 'Secondary'" label="Masters" name="masters">
          <UInputTags v-model="state.masters" icon="i-lucide-server" class="w-full" placeholder="Ex: 192.168.1.10, 192.168.1.11" />
        </UFormField>

        <UFormField v-if="state.kind === 'Primary'" label="Also Notify" name="also_notify">
          <UInputTags v-model="state.also_notify" icon="i-lucide-server" class="w-full" placeholder="Ex: 10.0.0.5, 10.0.0.6" />
        </UFormField>
        <UFormField label="SOA Edit API" name="soa_edit_api">
          <div class="flex items-center gap-2">
            <USelect v-model="state.soa_edit_api" icon="i-lucide-server-cog" :items="soaEditApis" class="w-full" />
            <UPopover arrow>
              <UButton icon="i-lucide-info" color="neutral" variant="subtle" class="rounded-full" />

              <template #content>
                <p class="max-w-xs text-sm break-words">
                  {{ soaDetails }}
                </p>
              </template>
            </UPopover>
          </div>
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
      <UButton label="Confirm" :loading="isLoading" @click="createZone" />
    </template>
  </UModal>
</template>
