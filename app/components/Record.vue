<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import type { Row } from '@tanstack/vue-table'
import { getPaginationRowModel } from '@tanstack/vue-table'

const toast = useToast()

const model = defineModel<string>('zoneId', { default: '' })

const UButton = resolveComponent('UButton')
const UDropdownMenu = resolveComponent('UDropdownMenu')
const UTooltip = resolveComponent('UTooltip')

const data: RecordForm[] = [
  { name: 'www', type: 'A', vl: '192.0.2.1', ttl: 3600, priority: 10, comment: 'Example record', zone: 'teste' },
  { name: 'mail', type: 'MX', vl: 'mail.example.com', ttl: 7200, priority: 5, comment: '', zone: 'teste' },
  { name: 'blog', type: 'CNAME', vl: 'blogs.example.com', ttl: 3600, priority: undefined, comment: 'Blog subdomain', zone: 'teste' },
]

const globalFilter = ref('')

const table = useTemplateRef('table')

const pagination = ref({ pageIndex: 0, pageSize: 10 })

watch(globalFilter, () => {
  pagination.value.pageIndex = 0
})

const columns: TableColumn<RecordForm>[] = [
  {
    accessorKey: 'name',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Name', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'type',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Type', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'value',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Value', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'ttl',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'TTL', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'priority',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Priority', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    id: 'age',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Age', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
  },
  {
    accessorKey: 'comment',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, { color: 'neutral', variant: 'ghost', label: 'Comment', icon: isSorted ? isSorted === 'asc' ? 'i-heroicons-bars-arrow-up' : 'i-heroicons-bars-arrow-down' : 'i-heroicons-arrows-up-down', class: '-mx-2.5', onClick: () => column.toggleSorting(column.getIsSorted() === 'asc') })
    },
    cell: ({ row }) => h(UTooltip, { text: row.original.comment || 'No comment', delayDuration: 0 }, () => h(UButton, { icon: 'i-lucide-eye', color: 'neutral', variant: 'ghost' })),
  },
  {
    id: 'actions',
    cell: ({ row }) => h(UDropdownMenu, { 'content': { align: 'end' }, 'items': getRowItems(row), 'aria-label': 'Actions dropdown' }, () => h(UButton, { 'icon': 'i-lucide-ellipsis-vertical', 'color': 'neutral', 'variant': 'ghost', 'class': 'ml-auto', 'aria-label': 'Actions dropdown' })),
  },
]

const { copy } = useClipboard()

function getRowItems(row: Row<RecordForm>){
  return [
    { type: 'label', label: `${row.original.name} actions` },
    {
      label: 'Copy Record Value',
      icon: 'i-lucide-copy',
      onSelect(){
        copy(row.original.vl || '')
        toast.add({ title: 'Record Value copied to clipboard!', color: 'success', icon: 'i-lucide-circle-check' })
      },
    },
    { type: 'separator' },
    { label: 'Edit Record', icon: 'i-lucide-pencil' },
    { label: 'Delete Record', icon: 'i-lucide-trash', color: 'error' },
  ]
}

const recordsOpts = ['A', 'AAAA', 'ALIAS', 'CAA', 'CNAME', 'HTTPS', 'MX', 'NS', 'TXT', 'SRV']

const state = ref<RecordForm>({ zone: '', name: '', type: 'A', vl: '', ttl: 60, priority: undefined, svcPriority: undefined, targetName: '', comment: '', port: undefined, weight: undefined, target: '', svcParams: '' })

const placeholder = computed(() => {
  switch (state.value.type){
    case 'A':
      return '192.0.2.1'
    case 'AAAA':
      return '2001:0db8:85a3:0000:0000:8a2e:0370:7334'
    case 'CNAME':
    case 'ALIAS':
      return 'host.example.com'
    case 'CAA':
      return '0 issue "letsencrypt.org"'
    case 'MX':
      return 'mail.example.com'
    case 'NS':
      return 'ns1.example.com'
    case 'TXT':
      return 'my example text record'
    default:
      return ''
  }
})
</script>

<template>
  <div class="p-4">
    <UButton variant="outline" class="mb-4" icon="i-lucide-arrow-left" label="Back to Zones" @click="model = ''" />

    <div>
      <h1 class="text-3xl font-bold">
        {{ model }} Records
      </h1>
      <p class="text-sm text-gray-500">
        Manage DNS records for this zone
      </p>
    </div>

    <UForm :schema="RecordSchema" :state="state" class="mt-6 space-y-4 p-5  rounded-lg bg-slate-950/40">
      <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
        <UFormField label="Name" name="name">
          <UInput v-model="state.name" icon="i-lucide-computer" class="w-full " placeholder="subdomain" />
        </UFormField>
        <UFormField label="Type" name="type">
          <USelect v-model="state.type" :items="recordsOpts" class="w-full" />
        </UFormField>
        <UFormField label="Value" name="vl">
          <UInput v-model="state.vl" icon="i-lucide-database" class="w-full" :placeholder="placeholder" :disabled="state.type === 'HTTPS' || state.type === 'SRV'" />
        </UFormField>
        <UFormField label="TTL" name="ttl">
          <UInputNumber v-model="state.ttl" :min="60" />
        </UFormField>
        <UFormField v-if="state.type !== 'HTTPS'" label="Priority" name="priority">
          <UInputNumber v-model="state.priority" :min="0" :disabled="state.type !== 'SRV' && state.type !== 'MX'" placeholder="10" />
        </UFormField>
      </div>

      <div v-if="state.type === 'HTTPS'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <UFormField label="SvcPriority" name="svcPriority">
          <UInputNumber v-model="state.svcPriority" :min="1" :max="65535" placeholder="0" class="w-full" />
        </UFormField>
        <UFormField label="TargetName" name="targetName">
          <UInput v-model="state.targetName" icon="i-lucide-target" class="w-full" placeholder="." />
        </UFormField>
        <UFormField label="SvcParams (Optional)" name="svcParams">
          <UInput v-model="state.svcParams" icon="i-lucide-settings" class="w-full" placeholder="alpn=h2,h3 foo=..." />
        </UFormField>
      </div>

      <div v-if="state.type === 'SRV'" class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <UFormField label="Weight" name="weight">
          <UInputNumber v-model="state.weight" :min="0" :max="65535" placeholder="10" class="w-full" />
        </UFormField>
        <UFormField label="Port" name="port">
          <UInputNumber v-model="state.port" :min="1" :max="65535" placeholder="80" class="w-full" />
        </UFormField>
        <UFormField label="Target" name="target">
          <UInput v-model="state.target" icon="i-lucide-target" class="w-full" placeholder="service.example.com" />
        </UFormField>
      </div>

      <UFormField label="Comment" name="comment">
        <UTextarea v-model="state.comment" icon="i-lucide-comment" class="w-full" placeholder="Optional comment" />
      </UFormField>

      <UButton variant="outline" class="mt-5 w-full flex justify-center" icon="i-lucide-plus" label="Add Record" />
    </UForm>

    <UTable ref="table" v-model:global-filter="globalFilter" v-model:pagination="pagination" :pagination-options="{ getPaginationRowModel: getPaginationRowModel()}" :data="data" :columns="columns" />

    <div v-if="data && data.length > pagination.pageSize" class="border-default flex justify-center border-t pt-4">
      <UPagination :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1" :items-per-page="table?.tableApi?.getState().pagination.pageSize" :total="table?.tableApi?.getFilteredRowModel().rows.length" @update:page="(p) => table?.tableApi?.setPageIndex(p - 1)" />
    </div>
  </div>
</template>
