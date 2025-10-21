<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import type { Row } from '@tanstack/vue-table'
import { getPaginationRowModel } from '@tanstack/vue-table'

const toast = useToast()
const { optionSelected } = await useConnection()
const { isLoading, start, finish } = useLoadingIndicator()

const model = defineModel<string>('zoneId', { default: '' })

const { data, refresh } = await useFetch<{ record: RecordForm[], soa: EditSOASchemaType }>('/server/api/zone/records', { method: 'GET', query: { connection: optionSelected, zone: model } })

const UButton = resolveComponent('UButton')
const UDropdownMenu = resolveComponent('UDropdownMenu')
const UTooltip = resolveComponent('UTooltip')

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
    accessorKey: 'vl',
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
      return '2001:0db8:85a3:00:00:0000:8a2e:0370:7334'
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

async function addRecord(){
  start()

  state.value.zone = model.value

  const body = RecordSchema.safeParse(state.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/zone/records', { method: 'PUT', body: body.data, query: { connection: optionSelected.value } })
    .catch(error => { toast.add({ title: error?.data?.message || error?.message || 'Error updating SOA record', icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  await refresh()
  state.value = { zone: '', name: '', type: 'A', vl: '', ttl: 60, priority: undefined, svcPriority: undefined, targetName: '', comment: '', port: undefined, weight: undefined, target: '', svcParams: '' }
  finish()
}

const modalEditSOA = ref(false)

const stateSOA = ref<EditSOASchemaType>({
  startOfAuthority: data.value?.soa.startOfAuthority || '',
  email: data.value?.soa.email || '',
  refresh: data.value?.soa.refresh || 0,
  retry: data.value?.soa.retry || 0,
  expire: data.value?.soa.expire || 0,
  negativeCacheTtl: data.value?.soa.negativeCacheTtl || 0,
})

async function updateSOA(){
  start()

  const body = EditSOASchema.safeParse(stateSOA.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/zone/soa', { method: 'PATCH', body: body.data, query: { connection: optionSelected.value, zone: model.value } })
    .catch(error => { toast.add({ title: error?.data?.message || error?.message || 'Error updating SOA record', icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  await refresh()
  modalEditSOA.value = false
  finish()
}

watch(modalEditSOA, nv => {
  if(nv) stateSOA.value = { startOfAuthority: data.value?.soa.startOfAuthority || '', email: data.value?.soa.email || '', refresh: data.value?.soa.refresh || 0, retry: data.value?.soa.retry || 0, expire: data.value?.soa.expire || 0, negativeCacheTtl: data.value?.soa.negativeCacheTtl || 0 }
})
</script>

<template>
  <UButton variant="outline" class="mb-4" icon="i-lucide-arrow-left" label="Back to Zones" @click="model = ''" />

  <div class="flex items-center justify-between my-6">
    <div>
      <h1 class="text-3xl font-bold">
        {{ model }} Records
      </h1>
      <p class="text-sm text-gray-500">
        Manage DNS records for this zone
      </p>
    </div>
    <UButton variant="outline" icon="i-lucide-pen" label="Edit SOA" :loading="isLoading" @click="modalEditSOA = true" />
  </div>

  <UForm :schema="RecordSchema" :state="state" class="mt-6 space-y-4 p-5 rounded-lg bg-slate-950/40">
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
      <UTextarea v-model="state.comment" class="w-full" placeholder="Optional comment" />
    </UFormField>

    <UButton variant="outline" class="mt-5 w-full flex justify-center" icon="i-lucide-plus" label="Add Record" :loading="isLoading" @click="addRecord" />
  </UForm>

  <UTable ref="table" v-model:global-filter="globalFilter" v-model:pagination="pagination" :pagination-options="{ getPaginationRowModel: getPaginationRowModel()}" :data="data?.record" :columns="columns" />

  <div v-if="data?.record && data.record.length > pagination.pageSize" class="border-default flex justify-center border-t pt-4">
    <UPagination :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1" :items-per-page="table?.tableApi?.getState().pagination.pageSize" :total="table?.tableApi?.getFilteredRowModel().rows.length" @update:page="(p) => table?.tableApi?.setPageIndex(p - 1)" />
  </div>

  <UModal v-model:open="modalEditSOA" title="Edit SOA" :description="`Edit the SOA record for ${model}`" :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :schema="EditSOASchema" :state="stateSOA" class="space-y-4">
        <UFormField label="Start of Authority" name="startOfAuthority">
          <UInput v-model="stateSOA.startOfAuthority" icon="i-lucide-shield-check" class="w-full" placeholder="Ex: ns1.example.com" />
        </UFormField>
        <UFormField label="Email" name="email">
          <UInput v-model="stateSOA.email" icon="i-lucide-mail" class="w-full" placeholder="Ex: hostmaster.example.com" />
        </UFormField>
        <UFormField label="Refresh" name="refresh">
          <UInputNumber v-model="stateSOA.refresh" :min="0" icon="i-lucide-refresh-cw" class="w-full" placeholder="3600" />
        </UFormField>
        <UFormField label="Retry" name="retry">
          <UInputNumber v-model="stateSOA.retry" :min="0" icon="i-lucide-clock" class="w-full" placeholder="600" />
        </UFormField>
        <UFormField label="Expire" name="expire">
          <UInputNumber v-model="stateSOA.expire" :min="0" icon="i-lucide-hourglass" class="w-full" placeholder="604800" />
        </UFormField>
        <UFormField label="Negative Cache TTL" name="negativeCacheTtl">
          <UInputNumber v-model="stateSOA.negativeCacheTtl" :min="0" icon="i-lucide-timer" class="w-full" placeholder="3600" />
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modalEditSOA = false" />
      <UButton label="Confirm" :loading="isLoading" @click="updateSOA" />
    </template>
  </UModal>
</template>
