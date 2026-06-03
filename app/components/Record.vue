<script setup lang="ts">
  import type { DropdownMenuItem, TableColumn } from '@nuxt/ui'
  import type { Row } from '@tanstack/vue-table'
  import { getPaginationRowModel } from '@tanstack/vue-table'
  import { safeParse } from 'valibot'

  import { UButton, UDropdownMenu, UPopover } from '#components'

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()

  const zoneId = defineModel<string>('zoneId', { required: true })
  const nivel = defineModel<string>('nivel', { required: true })

  const { data, refresh } = await useFetch<{ record: RecordForm[]; soa: EditSOASchemaType }>(
    '/server/api/records',
    { method: 'GET', query: { zone: zoneId } },
  )

  const globalFilter = ref('')

  const pagination = ref({ pageIndex: 0, pageSize: 20 })

  watch(globalFilter, () => {
    pagination.value.pageIndex = 0
  })

  const { copy } = useClipboard()

  const recordsOpts = [
    'A',
    'AAAA',
    'ALIAS',
    'CAA',
    'CNAME',
    'HTTPS',
    'MX',
    'NS',
    'PTR',
    'TXT',
    'SRV',
    'TLSA',
  ]

  const state = ref<RecordForm>({
    zone: '',
    name: '',
    type: 'A',
    vl: '',
    ttl: 3600,
    priority: undefined,
    svcPriority: undefined,
    targetName: '',
    comment: '',
    port: undefined,
    weight: undefined,
    target: '',
    svcParams: '',
  })

  const isEditing = ref(false)
  const oldState = ref<RecordForm>({
    zone: '',
    name: '',
    type: 'A',
    vl: '',
    ttl: 3600,
    priority: undefined,
    svcPriority: undefined,
    targetName: '',
    comment: '',
    port: undefined,
    weight: undefined,
    target: '',
    svcParams: '',
  })

  function cancelEdit(): void {
    isEditing.value = false
    state.value = {
      zone: '',
      name: '',
      type: 'A',
      vl: '',
      ttl: 3600,
      priority: undefined,
      svcPriority: undefined,
      targetName: '',
      comment: '',
      port: undefined,
      weight: undefined,
      target: '',
      svcParams: '',
    }
    oldState.value = {
      zone: '',
      name: '',
      type: 'A',
      vl: '',
      ttl: 3600,
      priority: undefined,
      svcPriority: undefined,
      targetName: '',
      comment: '',
      port: undefined,
      weight: undefined,
      target: '',
      svcParams: '',
    }
  }

  const placeholder = computed(() => {
    switch (state.value.type) {
      case 'A': {
        return '192.0.2.1'
      }
      case 'AAAA': {
        return '2001:0db8:85a3:00:00:0000:8a2e:0370:7334'
      }
      case 'CNAME':
      case 'ALIAS': {
        return 'host.example.com'
      }
      case 'CAA': {
        return '0 issue "letsencrypt.org"'
      }
      case 'MX': {
        return 'mail.example.com'
      }
      case 'NS': {
        return 'ns1.example.com'
      }
      case 'PTR': {
        return 'ptr.example.com'
      }
      case 'TXT': {
        return 'my example text record'
      }
      case 'TLSA': {
        return '3 1 1 a87f1b687ac...'
      }
      default: {
        return ''
      }
    }
  })

  async function addRecord(): Promise<void> {
    start()

    state.value.zone = zoneId.value

    if (!state.value.name || state.value.name === '') {
      state.value.name = zoneId.value
    } else if (state.value.name === zoneId.value) {
      state.value.name = zoneId.value
    } else if (!state.value.name.endsWith(`.${zoneId.value}`)) {
      state.value.name = `${state.value.name}.${zoneId.value}`
    }

    const body = safeParse(RecordSchema, state.value)
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/records', {
      method: 'PUT',
      body: body.output,
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao adicionar record',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    await refresh()
    state.value = {
      zone: '',
      name: '',
      type: 'A',
      vl: '',
      ttl: 3600,
      priority: undefined,
      svcPriority: undefined,
      targetName: '',
      comment: '',
      port: undefined,
      weight: undefined,
      target: '',
      svcParams: '',
    }
    finish()
  }

  async function editRecord(): Promise<void> {
    start()

    state.value.zone = zoneId.value
    oldState.value.zone = zoneId.value

    if (!state.value.name || state.value.name === '') {
      state.value.name = zoneId.value
    } else if (state.value.name === zoneId.value) {
      state.value.name = zoneId.value
    } else if (!state.value.name.endsWith(`.${zoneId.value}`)) {
      state.value.name = `${state.value.name}.${zoneId.value}`
    }

    if (!oldState.value.name || oldState.value.name === '') {
      oldState.value.name = zoneId.value
    } else if (!oldState.value.name.endsWith(`.${zoneId.value}`)) {
      oldState.value.name = `${oldState.value.name}.${zoneId.value}`
    }

    const body = safeParse(EditRecordSchema, { oldValue: oldState.value, newValue: state.value })

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/records', {
      method: 'PATCH',
      body: body.output,
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao atualizar record',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    await refresh()
    state.value = {
      zone: '',
      name: '',
      type: 'A',
      vl: '',
      ttl: 3600,
      priority: undefined,
      svcPriority: undefined,
      targetName: '',
      comment: '',
      port: undefined,
      weight: undefined,
      target: '',
      svcParams: '',
    }
    oldState.value = {
      zone: '',
      name: '',
      type: 'A',
      vl: '',
      ttl: 3600,
      priority: undefined,
      svcPriority: undefined,
      targetName: '',
      comment: '',
      port: undefined,
      weight: undefined,
      target: '',
      svcParams: '',
    }
    isEditing.value = false
    finish()
  }

  const modalEditSOA = ref(false)

  const stateSOA = ref<EditSOASchemaType>({
    startOfAuthority: data.value?.soa?.startOfAuthority || '',
    email: data.value?.soa?.email || '',
    refresh: data.value?.soa?.refresh || 0,
    retry: data.value?.soa?.retry || 0,
    expire: data.value?.soa?.expire || 0,
    negativeCacheTtl: data.value?.soa?.negativeCacheTtl || 0,
  })

  async function updateSOA(): Promise<void> {
    start()

    const body = safeParse(EditSOASchema, stateSOA.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/soa', {
      method: 'PATCH',
      body: body.output,
      query: { zone: zoneId.value },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao atualizar record',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    await refresh()
    modalEditSOA.value = false
    finish()
  }

  watch(modalEditSOA, (nv) => {
    if (nv) {
      stateSOA.value = {
        startOfAuthority: data.value?.soa?.startOfAuthority || '',
        email: data.value?.soa?.email || '',
        refresh: data.value?.soa?.refresh || 0,
        retry: data.value?.soa?.retry || 0,
        expire: data.value?.soa?.expire || 0,
        negativeCacheTtl: data.value?.soa?.negativeCacheTtl || 0,
      }
    }
  })

  const modalDelete = ref(false)
  const confirmDelete = ref('')
  const stateDelete = ref<RecordForm>({
    zone: '',
    name: '',
    type: 'A',
    vl: '',
    ttl: 3600,
    priority: undefined,
    svcPriority: undefined,
    targetName: '',
    comment: '',
    port: undefined,
    weight: undefined,
    target: '',
    svcParams: '',
  })

  function openDeleteModal(record: RecordForm): void {
    stateDelete.value = record
    modalDelete.value = true
  }

  async function removeRecord(): Promise<void> {
    start()

    if (confirmDelete.value !== stateDelete.value.name) {
      toast.add({
        title: 'Nome do registro não confere',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
      return finish({ error: true })
    }

    stateDelete.value.zone = zoneId.value

    const body = safeParse(RecordSchema, stateDelete.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/records', {
      method: 'DELETE',
      body: body.output,
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao remover record',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    await refresh()
    finish()
    modalDelete.value = false
  }

  watch(modalDelete, (nv) => {
    if (!nv) {
      stateDelete.value = {
        zone: '',
        name: '',
        type: 'A',
        vl: '',
        ttl: 3600,
        priority: undefined,
        svcPriority: undefined,
        targetName: '',
        comment: '',
        port: undefined,
        weight: undefined,
        target: '',
        svcParams: '',
      }
      confirmDelete.value = ''
    }
  })

  watch(isEditing, (nv) => {
    if (!nv) {
      state.value = {
        zone: '',
        name: '',
        type: 'A',
        vl: '',
        ttl: 3600,
        priority: undefined,
        svcPriority: undefined,
        targetName: '',
        comment: '',
        port: undefined,
        weight: undefined,
        target: '',
        svcParams: '',
      }
      oldState.value = {
        zone: '',
        name: '',
        type: 'A',
        vl: '',
        ttl: 3600,
        priority: undefined,
        svcPriority: undefined,
        targetName: '',
        comment: '',
        port: undefined,
        weight: undefined,
        target: '',
        svcParams: '',
      }
    }
  })

  const modalUsers = ref(false)

  const table = useTemplateRef('table')

  const columns: TableColumn<RecordForm>[] = [
    {
      accessorKey: 'name',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Nome',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
      cell: ({ row }) =>
        row.original.name?.split(zoneId.value)[0] === ''
          ? '@'
          : row.original.name?.split(`.${zoneId.value}`)[0],
    },
    {
      accessorKey: 'type',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Tipo',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
    },
    {
      accessorKey: 'vl',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Valor',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
      cell: ({ row }) => {
        const valor = row.original.vl || ''
        const isLong = valor.length > 30
        const display = isLong ? `${valor.slice(0, 30)}...` : valor

        if (!isLong) return h('p', { class: 'w-60 truncate' }, display)

        return h(
          UPopover,
          {
            mode: 'hover',
            ui: {
              content:
                'max-w-[320px] whitespace-pre-wrap break-all text-sm p-3 rounded-lg shadow-lg',
            },
          },
          {
            default: () =>
              h('p', { class: 'cursor-help dark:text-gray-400 w-60 truncate' }, display),
            content: () =>
              h(
                'p',
                { class: 'whitespace-pre-wrap break-all text-sm leading-snug dark:text-gray-200' },
                valor,
              ),
          },
        )
      },
    },
    {
      accessorKey: 'ttl',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'TTL',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
    },
    {
      accessorKey: 'comment',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Descrição',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
      cell: ({ row }) => {
        const comment = row.original.comment || ''
        const isLong = comment.length > 30
        const display = isLong ? `${comment.slice(0, 30)}...` : comment

        if (!isLong) return h('p', { class: 'w-60 truncate' }, display)

        return h(
          UPopover,
          {
            mode: 'hover',
            ui: {
              content:
                'max-w-[320px] whitespace-pre-wrap break-all text-sm p-3 rounded-lg shadow-lg',
            },
          },
          {
            default: () =>
              h('p', { class: 'cursor-help dark:text-gray-400 w-60 truncate' }, display),
            content: () =>
              h(
                'p',
                { class: 'whitespace-pre-wrap break-all text-sm leading-snug dark:text-gray-200' },
                comment,
              ),
          },
        )
      },
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        if (nivel.value !== 'LEITURA') {
          return h(
            UDropdownMenu,
            {
              content: { align: 'end' },
              items: getRowItems(row),
              'aria-label': 'Actions dropdown',
            },
            () =>
              h(UButton, {
                icon: 'i-lucide-ellipsis-vertical',
                color: 'neutral',
                variant: 'ghost',
                class: 'ml-auto',
                'aria-label': 'Actions dropdown',
              }),
          )
        }
        return null
      },
    },
  ]

  function getRowItems(row: Row<RecordForm>): DropdownMenuItem[] {
    return [
      { type: 'label', label: `Ações` },
      {
        label: 'Copiar valor do record',
        icon: 'i-lucide-copy',
        onSelect(): void {
          copy(row.original.vl || '')
          toast.add({
            title: 'Valor do record copiado para a área de transferência!',
            color: 'success',
            icon: 'i-lucide-circle-check',
          })
        },
      },
      { type: 'separator' },
      {
        label: 'Editar Record',
        icon: 'i-lucide-pencil',
        onSelect(): void {
          isEditing.value = true
          let oldMx = ''
          let newMx = ''
          if (row.original.type === 'MX' && row.original.vl) {
            oldMx = row.original.vl.split(' ')[1] || ''
            newMx = oldMx
          }
          state.value = {
            ...row.original,
            vl: row.original.type === 'MX' && row.original.vl ? newMx : row.original.vl,
            name:
              row.original.name === zoneId.value
                ? ''
                : row.original.name?.split(`.${zoneId.value}`)[0] || '',
          }
          oldState.value = {
            ...row.original,
            vl: row.original.type === 'MX' && row.original.vl ? oldMx : row.original.vl,
            name:
              row.original.name === zoneId.value
                ? ''
                : row.original.name?.split(`.${zoneId.value}`)[0] || '',
          }
        },
      },
      {
        label: 'Delete Record',
        icon: 'i-lucide-trash',
        color: 'error',
        onSelect(): void {
          openDeleteModal(row.original)
        },
      },
    ]
  }
</script>

<template>
  <UButton
    variant="outline"
    class="mb-4"
    icon="i-lucide-arrow-left"
    label="Voltar para Zonas"
    @click="zoneId = ''" />

  <div class="my-6 flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-bold">{{ zoneId }} Records</h1>
      <p class="text-sm text-gray-500">
        Controle e gerencie os records DNS para o seu domínio aqui.
      </p>
    </div>
    <div class="flex flex-wrap items-center justify-center gap-2">
      <UButton
        v-if="nivel === 'ADMINISTRADOR'"
        variant="outline"
        icon="i-lucide-users"
        label="Usuários"
        :loading="isLoading"
        @click="modalUsers = true" />
      <UButton
        v-if="nivel === 'ADMINISTRADOR'"
        variant="outline"
        icon="i-lucide-pen"
        label="Editar SOA"
        :loading="isLoading"
        @click="modalEditSOA = true" />
    </div>
  </div>

  <UForm
    v-if="nivel !== 'LEITURA'"
    :schema="RecordSchema"
    :state="state"
    class="mt-6 space-y-4 rounded-lg bg-slate-100 p-5 dark:bg-slate-950/40"
    @submit="isEditing ? (modalEditSOA = true) : addRecord">
    <div class="grid grid-cols-1 gap-4 md:grid-cols-3 lg:grid-cols-5">
      <UFormField label="Nome" name="name">
        <UInput
          v-model="state.name"
          :disabled="isEditing"
          icon="i-lucide-computer"
          class="w-full"
          placeholder="subdomínio" />
      </UFormField>
      <UFormField label="Tipo" name="type">
        <USelect v-model="state.type" :disabled="isEditing" :items="recordsOpts" class="w-full" />
      </UFormField>
      <UFormField label="Valor" name="vl">
        <UInput
          v-model="state.vl"
          icon="i-lucide-database"
          class="w-full"
          :placeholder="placeholder"
          :disabled="state.type === 'HTTPS' || state.type === 'SRV'" />
      </UFormField>
      <UFormField label="TTL" name="ttl">
        <UInputNumber v-model="state.ttl" :min="60" />
      </UFormField>
      <UFormField v-if="state.type !== 'HTTPS'" label="Prioridade" name="priority">
        <UInputNumber
          v-model="state.priority"
          :min="0"
          :disabled="state.type !== 'SRV' && state.type !== 'MX'"
          placeholder="10" />
      </UFormField>
    </div>

    <div v-if="state.type === 'HTTPS'" class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
      <UFormField label="SvcPriority" name="svcPriority">
        <UInputNumber
          v-model="state.svcPriority"
          :min="1"
          :max="65535"
          placeholder="0"
          class="w-full" />
      </UFormField>
      <UFormField label="TargetName" name="targetName">
        <UInput v-model="state.targetName" icon="i-lucide-target" class="w-full" placeholder="." />
      </UFormField>
      <UFormField label="SvcParams (Opcional)" name="svcParams">
        <UInput
          v-model="state.svcParams"
          icon="i-lucide-settings"
          class="w-full"
          placeholder="alpn=h2,h3 foo=..." />
      </UFormField>
    </div>

    <div v-if="state.type === 'SRV'" class="grid grid-cols-1 gap-4 md:grid-cols-3">
      <UFormField label="Weight" name="weight">
        <UInputNumber
          v-model="state.weight"
          :min="0"
          :max="65535"
          placeholder="10"
          class="w-full" />
      </UFormField>
      <UFormField label="Port" name="port">
        <UInputNumber v-model="state.port" :min="1" :max="65535" placeholder="80" class="w-full" />
      </UFormField>
      <UFormField label="Target" name="target">
        <UInput
          v-model="state.target"
          icon="i-lucide-target"
          class="w-full"
          placeholder="service.example.com" />
      </UFormField>
    </div>

    <UFormField label="Comentário" name="comment">
      <UTextarea v-model="state.comment" class="w-full" placeholder="Comentário opcional" />
    </UFormField>

    <div v-if="isEditing" class="flex w-full items-center gap-2">
      <UButton
        variant="outline"
        class="mt-5 flex w-full justify-center"
        icon="i-lucide-pen"
        label="Editar Record"
        :loading="isLoading"
        @click="editRecord" />
      <UButton
        variant="outline"
        color="error"
        class="mt-5 flex w-full max-w-32 justify-center"
        icon="i-lucide-x"
        label="Cancelar"
        :loading="isLoading"
        @click="cancelEdit" />
    </div>
    <UButton
      v-else
      variant="outline"
      class="mt-5 flex w-full justify-center"
      icon="i-lucide-plus"
      label="Adicionar Record"
      :loading="isLoading"
      @click="addRecord" />
  </UForm>

  <div class="space-x-6">
    <UInput
      v-model="globalFilter"
      class="mt-10 mb-4"
      placeholder="Buscar records..."
      icon="i-lucide-search" />
  </div>

  <ClientOnly>
    <UTable
      ref="table"
      v-model:global-filter="globalFilter"
      v-model:pagination="pagination"
      class="mb-10"
      :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
      :data="data?.record"
      :columns="columns" />
  </ClientOnly>

  <div
    v-if="data?.record && data.record.length > pagination.pageSize"
    class="flex justify-center border-t border-default pt-4">
    <UPagination
      active-
      active-variant="subtle"
      :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1"
      :items-per-page="table?.tableApi?.getState().pagination.pageSize"
      :total="table?.tableApi?.getFilteredRowModel().rows.length"
      @update:page="(p) => table?.tableApi?.setPageIndex(p - 1)" />
  </div>

  <UModal
    v-model:open="modalEditSOA"
    title="Editar SOA"
    :description="`Editar o registro SOA para ${zoneId}`"
    :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :schema="EditSOASchema" :state="stateSOA" class="space-y-4">
        <UFormField label="Start of Authority" name="startOfAuthority">
          <UInput
            v-model="stateSOA.startOfAuthority"
            icon="i-lucide-shield-check"
            class="w-full"
            placeholder="Ex: ns1.example.com" />
        </UFormField>
        <UFormField label="Email" name="email">
          <UInput
            v-model="stateSOA.email"
            icon="i-lucide-mail"
            class="w-full"
            placeholder="Ex: hostmaster.example.com" />
        </UFormField>
        <UFormField label="Refresh" name="refresh">
          <UInputNumber
            v-model="stateSOA.refresh"
            :min="0"
            icon="i-lucide-refresh-cw"
            class="w-full"
            placeholder="3600" />
        </UFormField>
        <UFormField label="Retry" name="retry">
          <UInputNumber
            v-model="stateSOA.retry"
            :min="0"
            icon="i-lucide-clock"
            class="w-full"
            placeholder="600" />
        </UFormField>
        <UFormField label="Expire" name="expire">
          <UInputNumber
            v-model="stateSOA.expire"
            :min="0"
            icon="i-lucide-hourglass"
            class="w-full"
            placeholder="604800" />
        </UFormField>
        <UFormField label="Negative Cache TTL" name="negativeCacheTtl">
          <UInputNumber
            v-model="stateSOA.negativeCacheTtl"
            :min="0"
            icon="i-lucide-timer"
            class="w-full"
            placeholder="3600" />
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <UButton
        label="Cancel"
        :loading="isLoading"
        variant="outline"
        @click="modalEditSOA = false" />
      <UButton label="Confirm" :loading="isLoading" @click="updateSOA" />
    </template>
  </UModal>

  <UModal
    v-model:open="modalDelete"
    title="Aviso"
    description="Você está prestes a deletar um registro, esta ação não pode ser desfeita."
    :ui="{ footer: 'justify-end' }">
    <template #body>
      <p class="dark:text-gray-200">
        If you are sure you want to continue, write the name of your record below
        <span class="font-bold">'{{ stateDelete.name }}'</span>.
      </p>
      <UInput v-model="confirmDelete" class="mt-2 w-full" color="error" placeholder="Record Name" />
    </template>

    <template #footer>
      <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modalDelete = false" />
      <UButton
        label="Confirm"
        color="error"
        :loading="isLoading"
        :disabled="confirmDelete !== stateDelete.name"
        @click="removeRecord" />
    </template>
  </UModal>

  <Users v-if="nivel === 'ADMINISTRADOR'" v-model:open="modalUsers" v-model:zone-id="zoneId" />
</template>
