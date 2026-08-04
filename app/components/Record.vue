<script setup lang="ts">
  import type { DropdownMenuItem, TableColumn } from '@nuxt/ui'
  import type { Column, Row } from '@tanstack/vue-table'
  import { getPaginationRowModel } from '@tanstack/vue-table'
  import { safeParse } from 'valibot'
  import type { VNode } from 'vue'

  import { UButton, UDropdownMenu, UPopover } from '#components'

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()
  const { recordSchema, editSOASchema, editRecordSchema } = useRecordSchema()
  const { recordTypes, typeMeta } = useRecordTypes()

  const props = defineProps({
    zone: { type: String, required: true },
  })

  const zoneId = computed(() => props.zone)

  const { data, refresh } = await useApi<{
    record: RecordSchema[]
    soa: EditSOASchema
    nivel: string
  }>('/records', { method: 'GET', query: { zone: zoneId } })

  const nivel = computed(() => data.value?.nivel ?? '')

  const globalFilter = ref('')

  const pagination = ref({ pageIndex: 0, pageSize: 20 })

  watch(globalFilter, () => {
    pagination.value.pageIndex = 0
  })

  const { copy } = useClipboard()

  function blankRecord(): RecordSchema {
    return {
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

  const state = ref<RecordSchema>(blankRecord())

  const isEditing = ref(false)
  const oldState = ref<RecordSchema>(blankRecord())

  const formOpen = ref(false)

  function cancelEdit(): void {
    isEditing.value = false
    state.value = blankRecord()
    oldState.value = blankRecord()
  }

  function openCreate(): void {
    cancelEdit()
    formOpen.value = true
  }

  watch(formOpen, (nv) => {
    if (!nv) isEditing.value = false
  })

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

  function qualifyName(value: string | undefined): string {
    if (!value || value === '') return zoneId.value
    if (value === zoneId.value) return zoneId.value
    if (!value.endsWith(`.${zoneId.value}`)) return `${value}.${zoneId.value}`
    return value
  }

  async function addRecord(): Promise<void> {
    start()

    state.value.zone = zoneId.value
    state.value.name = qualifyName(state.value.name)

    const body = safeParse(recordSchema, state.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/records', {
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
    state.value = blankRecord()
    formOpen.value = false
    finish()
  }

  async function editRecord(): Promise<void> {
    start()

    state.value.zone = zoneId.value
    oldState.value.zone = zoneId.value

    state.value.name = qualifyName(state.value.name)
    oldState.value.name = qualifyName(oldState.value.name)

    const body = safeParse(editRecordSchema, { oldValue: oldState.value, newValue: state.value })

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/records', {
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
    formOpen.value = false
    isEditing.value = false
    finish()
  }

  const modalEditSOA = ref(false)

  function soaFromData(): EditSOASchema {
    return {
      startOfAuthority: data.value?.soa?.startOfAuthority || '',
      email: data.value?.soa?.email || '',
      refresh: data.value?.soa?.refresh || 0,
      retry: data.value?.soa?.retry || 0,
      expire: data.value?.soa?.expire || 0,
      negativeCacheTtl: data.value?.soa?.negativeCacheTtl || 0,
    }
  }

  const stateSOA = ref<EditSOASchema>(soaFromData())

  async function updateSOA(): Promise<void> {
    start()

    const body = safeParse(editSOASchema, stateSOA.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/soa', {
      method: 'PATCH',
      body: body.output,
      query: { zone: zoneId.value },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao atualizar SOA',
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
    if (nv) stateSOA.value = soaFromData()
  })

  const modalDelete = ref(false)
  const confirmDelete = ref('')
  const stateDelete = ref<RecordSchema>(blankRecord())

  function openDeleteModal(record: RecordSchema): void {
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

    const body = safeParse(recordSchema, stateDelete.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/records', {
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
      stateDelete.value = blankRecord()
      confirmDelete.value = ''
    }
  })

  watch(isEditing, (nv) => {
    if (!nv) {
      state.value = blankRecord()
      oldState.value = blankRecord()
    }
  })

  const modalUsers = ref(false)

  function openUsers(): void {
    modalUsers.value = true
  }

  function openSOA(): void {
    modalEditSOA.value = true
  }

  const modalReversos = ref(false)
  const reversosAusentes = ref<ReversoAusente[]>([])
  const reversosSelecionados = ref<string[]>([])

  function toggleReverso(nomeReverso: string, marcado: boolean): void {
    reversosSelecionados.value = marcado
      ? [...reversosSelecionados.value, nomeReverso]
      : reversosSelecionados.value.filter((item) => item !== nomeReverso)
  }

  const selecaoReversos = computed(() => {
    if (reversosSelecionados.value.length === 0) return false
    if (reversosSelecionados.value.length === reversosAusentes.value.length) return true
    return 'indeterminate' as const
  })

  function toggleTodosReversos(marcado: boolean | 'indeterminate'): void {
    reversosSelecionados.value =
      marcado === true ? reversosAusentes.value.map((reverso) => reverso.nomeReverso) : []
  }

  async function checkReversos(): Promise<void> {
    start()

    const res = await $api<ReversoAusente[]>('/reverses', {
      method: 'GET',
      query: { zone: zoneId.value },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao verificar reversos',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    reversosAusentes.value = res
    reversosSelecionados.value = []
    modalReversos.value = true
    finish()
  }

  async function createReversos(): Promise<void> {
    start()

    const res = await $api('/reverses', {
      method: 'PUT',
      query: { zone: zoneId.value },
      body: { reversos: reversosSelecionados.value },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao criar reversos',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    modalReversos.value = false
    reversosAusentes.value = []
    reversosSelecionados.value = []
    await refresh()
    finish()
  }

  const modalOrfaos = ref(false)
  const reversosOrfaos = ref<ReversoOrfao[]>([])

  async function checkOrfaos(): Promise<void> {
    start()

    const res = await $api<ReversoOrfao[]>('/reverses/orphans', {
      method: 'GET',
      query: { zone: zoneId.value },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao verificar reversos órfãos',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    reversosOrfaos.value = res
    modalOrfaos.value = true
    finish()
  }

  async function deleteOrfao(nome: string, close: () => void): Promise<void> {
    close()
    start()

    const res = await $api('/reverses', {
      method: 'DELETE',
      query: { zone: zoneId.value, name: nome },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao excluir reverso',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    reversosOrfaos.value = reversosOrfaos.value.filter((orfao) => orfao.nomeReverso !== nome)
    await refresh()
    finish()
  }

  async function backToZones(): Promise<void> {
    await navigateTo('/')
  }

  const table = useTemplateRef('table')

  function shortName(name: string | undefined): string {
    if (name?.split(zoneId.value)[0] === '') return '@'
    return name?.split(`.${zoneId.value}`)[0] || '@'
  }

  function sortableHeader(column: Column<RecordSchema>, label: string): VNode {
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

  function truncatedCell(value: string, extraClass = ''): VNode {
    const isLong = value.length > 30
    const display = isLong ? `${value.slice(0, 30)}…` : value

    if (!isLong) return h('p', { class: `w-60 truncate ${extraClass}` }, display)

    return h(
      UPopover,
      {
        mode: 'hover',
        ui: { content: 'max-w-[360px] p-3 rounded-lg shadow-lg' },
      },
      {
        default: () => h('p', { class: `w-60 cursor-help truncate ${extraClass}` }, display),
        content: () =>
          h(
            'p',
            { class: 'data text-xs leading-relaxed break-all whitespace-pre-wrap text-toned' },
            value,
          ),
      },
    )
  }

  const columns: TableColumn<RecordSchema>[] = [
    {
      accessorKey: 'name',
      header: ({ column }) => sortableHeader(column, 'Nome'),
      cell: ({ row }) =>
        h('span', { class: 'data font-medium text-highlighted' }, shortName(row.original.name)),
    },
    {
      accessorKey: 'type',
      header: ({ column }) => sortableHeader(column, 'Tipo'),
      cell: ({ row }) =>
        h(
          'span',
          {
            class: `data inline-flex items-center rounded-md px-2 py-0.5 text-[11px] font-semibold ring-1 ring-inset ${typeMeta(row.original.type).chip}`,
          },
          row.original.type,
        ),
    },
    {
      accessorKey: 'vl',
      header: ({ column }) => sortableHeader(column, 'Valor'),
      cell: ({ row }) => truncatedCell(row.original.vl || '', 'data text-toned'),
    },
    {
      accessorKey: 'ttl',
      header: ({ column }) => sortableHeader(column, 'TTL'),
      cell: ({ row }) => h('span', { class: 'data tnum text-muted' }, row.original.ttl),
    },
    {
      accessorKey: 'comment',
      header: ({ column }) => sortableHeader(column, 'Descrição'),
      cell: ({ row }) => truncatedCell(row.original.comment || '', 'text-dimmed'),
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        if (nivel.value === 'LEITURA') return null

        return h(
          'div',
          { class: 'text-right' },
          h(
            UDropdownMenu,
            {
              content: { align: 'end' },
              items: getRowItems(row),
              'aria-label': 'Ações do record',
            },
            () =>
              h(UButton, {
                icon: 'i-lucide-ellipsis-vertical',
                color: 'neutral',
                variant: 'ghost',
                size: 'sm',
                class: 'ml-auto',
                'aria-label': `Ações do record ${row.original.name}`,
              }),
          ),
        )
      },
    },
  ]

  function getRowItems(row: Row<RecordSchema>): DropdownMenuItem[] {
    return [
      { type: 'label', label: shortName(row.original.name) },
      {
        label: 'Copiar valor',
        icon: 'i-lucide-copy',
        onSelect(): void {
          copy(row.original.vl || '')
          toast.add({
            title: 'Valor copiado para a área de transferência',
            color: 'success',
            icon: 'i-lucide-circle-check',
          })
        },
      },
      { type: 'separator' },
      {
        label: 'Editar record',
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
          formOpen.value = true
        },
      },
      { type: 'separator' },
      {
        label: 'Deletar record',
        icon: 'i-lucide-trash-2',
        color: 'error',
        onSelect(): void {
          openDeleteModal(row.original)
        },
      },
    ]
  }

  const totalRecords = computed(() => data.value?.record?.length ?? 0)

  const isAdmin = computed(() => nivel.value === 'ADMINISTRADOR')
  const canWrite = computed(() => nivel.value !== 'LEITURA')
  const isReverse = computed(() => /\.(?:in-addr|ip6)\.arpa\.?$/u.test(zoneId.value))
  const isReverseIpv6 = computed(() => /\.ip6\.arpa\.?$/u.test(zoneId.value))

  const zoneLabel = computed(() => zoneId.value)
  const nameHelp = computed(() => qualifyName(state.value.name))
</script>

<template>
  <div class="space-y-5">
    <UButton
      variant="link"
      color="neutral"
      size="sm"
      class="-ml-2 gap-1.5"
      icon="i-lucide-arrow-left"
      label="Todas as zonas"
      @click="backToZones" />

    <PageHeader eyebrow="Zona" :title="zoneLabel" mono description="Registros DNS desta zona.">
      <template #actions>
        <UButton
          v-if="isAdmin"
          variant="outline"
          color="neutral"
          icon="i-lucide-users"
          label="Usuários"
          :loading="isLoading"
          @click="openUsers" />
        <UButton
          v-if="isAdmin"
          variant="outline"
          color="neutral"
          icon="i-lucide-file-cog"
          label="Editar SOA"
          :loading="isLoading"
          @click="openSOA" />
        <UButton
          v-if="canWrite && !isReverse"
          variant="outline"
          color="neutral"
          icon="i-lucide-network"
          label="Verificar reversos"
          :loading="isLoading"
          @click="checkReversos" />
        <UButton
          v-if="canWrite && isReverse"
          variant="outline"
          color="neutral"
          icon="i-lucide-unlink"
          label="Verificar órfãos"
          :loading="isLoading"
          @click="checkOrfaos" />
        <UButton
          v-if="canWrite"
          icon="i-lucide-plus"
          label="Novo record"
          :loading="isLoading"
          @click="openCreate" />
      </template>
    </PageHeader>

    <div class="overflow-hidden rounded-xl border border-default bg-default">
      <div
        class="flex flex-col gap-3 border-b border-default bg-muted/40 p-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-2.5">
          <UIcon name="i-lucide-list" class="size-4 shrink-0 text-dimmed" />
          <p class="text-sm font-medium text-toned">
            <span class="data text-highlighted tnum">{{ totalRecords }}</span>
            {{ totalRecords === 1 ? 'record' : 'records' }}
          </p>
        </div>

        <UInput
          v-model="globalFilter"
          icon="i-lucide-search"
          placeholder="Filtrar por nome, tipo ou valor..."
          class="w-full sm:w-80"
          :ui="{ base: 'data' }" />
      </div>

      <ClientOnly>
        <UTable
          ref="table"
          v-model:global-filter="globalFilter"
          v-model:pagination="pagination"
          :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
          :data="data?.record"
          :columns="columns"
          :ui="{
            tr: 'transition-colors duration-150 hover:bg-elevated/60',
            th: 'py-2',
            td: 'py-2.5',
          }">
          <template #empty>
            <div class="flex flex-col items-center gap-2 py-12 text-center">
              <UIcon name="i-lucide-file-search" class="size-7 text-dimmed" />
              <p class="text-sm font-medium text-toned">Nenhum record encontrado</p>
              <p class="max-w-xs text-xs text-dimmed">
                {{
                  globalFilter
                    ? 'Nenhum resultado para esse filtro. Tente outro termo.'
                    : 'Esta zona ainda não possui registros DNS.'
                }}
              </p>
              <UButton
                v-if="canWrite && !globalFilter"
                class="mt-2"
                size="sm"
                variant="outline"
                icon="i-lucide-plus"
                label="Adicionar o primeiro"
                @click="openCreate" />
            </div>
          </template>
        </UTable>
      </ClientOnly>

      <div
        v-if="data?.record && data.record.length > pagination.pageSize"
        class="flex justify-center border-t border-default bg-muted/40 p-3">
        <UPagination
          active-variant="subtle"
          :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1"
          :items-per-page="table?.tableApi?.getState().pagination.pageSize"
          :total="table?.tableApi?.getFilteredRowModel().rows.length"
          @update:page="(p) => table?.tableApi?.setPageIndex(p - 1)" />
      </div>
    </div>
  </div>

  <USlideover
    v-model:open="formOpen"
    :title="isEditing ? 'Editar record' : 'Novo record'"
    :description="`Zona ${zoneLabel}`"
    :ui="{ content: 'max-w-xl', footer: 'justify-end' }">
    <template #body>
      <UForm
        :schema="recordSchema"
        :state="state"
        class="space-y-5"
        @submit="isEditing ? editRecord() : addRecord()">
        <div class="grid gap-4 sm:grid-cols-2">
          <UFormField label="Nome" name="name" hint="subdomínio">
            <UInput
              v-model="state.name"
              :disabled="isEditing"
              icon="i-lucide-tag"
              class="w-full"
              placeholder="www"
              :ui="{ base: 'data' }" />
            <template v-if="!isReverseIpv6" #help>
              <span
                :title="qualifyName(state.name)"
                class="block truncate data text-[11px] text-dimmed">
                {{ nameHelp }}
              </span>
            </template>
          </UFormField>

          <UFormField label="Tipo" name="type">
            <USelect
              v-model="state.type"
              :disabled="isEditing"
              :items="recordTypes"
              class="w-full"
              :ui="{ base: 'data' }" />
            <template #help>
              <span class="text-[11px] text-dimmed">{{ typeMeta(state.type).hint }}</span>
            </template>
          </UFormField>
        </div>

        <UFormField
          v-if="state.type !== 'HTTPS' && state.type !== 'SRV'"
          label="Valor"
          name="vl"
          required>
          <UTextarea
            v-if="state.type === 'TXT'"
            v-model="state.vl"
            class="w-full"
            :rows="4"
            :placeholder="placeholder"
            :ui="{ base: 'data' }" />
          <UInput
            v-else
            v-model="state.vl"
            icon="i-lucide-database"
            class="w-full"
            :placeholder="placeholder"
            :ui="{ base: 'data' }" />
        </UFormField>

        <div v-if="state.type === 'HTTPS'" class="space-y-4 rounded-lg bg-muted/50 p-4">
          <div class="grid gap-4 sm:grid-cols-2">
            <UFormField label="SvcPriority" name="svcPriority" required>
              <UInputNumber
                v-model="state.svcPriority"
                :min="1"
                :max="65535"
                class="w-full"
                placeholder="1" />
            </UFormField>
            <UFormField label="TargetName" name="targetName" required>
              <UInput
                v-model="state.targetName"
                icon="i-lucide-target"
                class="w-full"
                placeholder="."
                :ui="{ base: 'data' }" />
            </UFormField>
          </div>
          <UFormField label="SvcParams" name="svcParams" hint="opcional">
            <UInput
              v-model="state.svcParams"
              icon="i-lucide-settings-2"
              class="w-full"
              placeholder="alpn=h2,h3"
              :ui="{ base: 'data' }" />
          </UFormField>
        </div>

        <div v-if="state.type === 'SRV'" class="space-y-4 rounded-lg bg-muted/50 p-4">
          <div class="grid gap-4 sm:grid-cols-2">
            <UFormField label="Weight" name="weight" required>
              <UInputNumber
                v-model="state.weight"
                :min="0"
                :max="65535"
                class="w-full"
                placeholder="10" />
            </UFormField>
            <UFormField label="Port" name="port" required>
              <UInputNumber
                v-model="state.port"
                :min="1"
                :max="65535"
                class="w-full"
                placeholder="443" />
            </UFormField>
          </div>
          <UFormField label="Target" name="target" required>
            <UInput
              v-model="state.target"
              icon="i-lucide-target"
              class="w-full"
              placeholder="service.example.com"
              :ui="{ base: 'data' }" />
          </UFormField>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <UFormField label="TTL" name="ttl" hint="segundos">
            <UInputNumber v-model="state.ttl" :min="60" class="w-full" />
          </UFormField>

          <UFormField
            v-if="state.type !== 'HTTPS'"
            label="Prioridade"
            name="priority"
            :hint="state.type === 'SRV' || state.type === 'MX' ? undefined : 'n/a'">
            <UInputNumber
              v-model="state.priority"
              :min="0"
              :disabled="state.type !== 'SRV' && state.type !== 'MX'"
              class="w-full"
              placeholder="10" />
          </UFormField>
        </div>

        <UFormField label="Comentário" name="comment" hint="opcional">
          <UTextarea
            v-model="state.comment"
            class="w-full"
            :rows="3"
            placeholder="Para que serve este record" />
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <UButton
        label="Cancelar"
        color="neutral"
        variant="ghost"
        :loading="isLoading"
        @click="
          () => {
            formOpen = false
          }
        " />
      <UButton
        :label="isEditing ? 'Salvar alterações' : 'Adicionar record'"
        :icon="isEditing ? 'i-lucide-check' : 'i-lucide-plus'"
        :loading="isLoading"
        @click="isEditing ? editRecord() : addRecord()" />
    </template>
  </USlideover>

  <UModal
    v-model:open="modalEditSOA"
    title="Editar SOA"
    :description="`Registro Start of Authority da zona ${zoneId}`"
    :ui="{ content: 'max-w-2xl', footer: 'justify-end' }">
    <template #body>
      <UForm :schema="editSOASchema" :state="stateSOA" class="grid gap-4 sm:grid-cols-2">
        <UFormField
          label="Servidor primário"
          name="startOfAuthority"
          class="sm:col-span-2"
          required>
          <UInput
            v-model="stateSOA.startOfAuthority"
            icon="i-lucide-server"
            class="w-full"
            placeholder="ns1.example.com"
            :ui="{ base: 'data' }" />
        </UFormField>

        <UFormField label="Email do responsável" name="email" class="sm:col-span-2" required>
          <UInput
            v-model="stateSOA.email"
            icon="i-lucide-mail"
            class="w-full"
            placeholder="hostmaster.example.com"
            :ui="{ base: 'data' }" />
        </UFormField>

        <UFormField label="Refresh" name="refresh" hint="segundos">
          <UInputNumber v-model="stateSOA.refresh" :min="0" class="w-full" placeholder="3600" />
        </UFormField>

        <UFormField label="Retry" name="retry" hint="segundos">
          <UInputNumber v-model="stateSOA.retry" :min="0" class="w-full" placeholder="600" />
        </UFormField>

        <UFormField label="Expire" name="expire" hint="segundos">
          <UInputNumber v-model="stateSOA.expire" :min="0" class="w-full" placeholder="604800" />
        </UFormField>

        <UFormField label="Negative Cache TTL" name="negativeCacheTtl" hint="segundos">
          <UInputNumber
            v-model="stateSOA.negativeCacheTtl"
            :min="0"
            class="w-full"
            placeholder="86400" />
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <UButton
        label="Cancelar"
        color="neutral"
        variant="ghost"
        :loading="isLoading"
        @click="
          () => {
            modalEditSOA = false
          }
        " />
      <UButton label="Salvar SOA" icon="i-lucide-check" :loading="isLoading" @click="updateSOA" />
    </template>
  </UModal>

  <UModal
    v-model:open="modalDelete"
    title="Deletar record"
    description="Esta ação é permanente e não pode ser desfeita."
    :ui="{ footer: 'justify-end' }">
    <template #body>
      <div class="space-y-4">
        <div class="rounded-lg border border-default bg-muted/50 p-3">
          <span
            class="inline-flex items-center rounded-md px-2 py-0.5 data text-[11px] font-semibold ring-1 ring-inset"
            :class="typeMeta(stateDelete.type).chip">
            {{ stateDelete.type }}
          </span>
          <p class="mt-2 data text-sm font-medium break-all text-highlighted">
            {{ stateDelete.name }}
          </p>
          <p v-if="stateDelete.vl" class="mt-2 data text-xs break-all text-muted">
            {{ stateDelete.vl }}
          </p>
        </div>

        <UFormField label="Confirme digitando o nome do record">
          <template #description>
            <p class="text-xs text-dimmed">
              Digite o nome do record para confirmar a exclusão. Esta ação não pode ser desfeita.
            </p>
          </template>
          <UInput
            v-model="confirmDelete"
            class="w-full"
            color="error"
            autocomplete="off"
            :placeholder="stateDelete.name"
            :ui="{ base: 'data' }" />
        </UFormField>
      </div>
    </template>

    <template #footer>
      <UButton
        label="Cancelar"
        color="neutral"
        variant="ghost"
        :loading="isLoading"
        @click="
          () => {
            modalDelete = false
          }
        " />
      <UButton
        label="Deletar record"
        icon="i-lucide-trash-2"
        color="error"
        :loading="isLoading"
        :disabled="confirmDelete !== stateDelete.name"
        @click="removeRecord" />
    </template>
  </UModal>

  <UModal
    v-model:open="modalReversos"
    title="Reversos ausentes"
    :description="`Registros A e AAAA de ${zoneLabel} sem PTR na zona reversa`"
    :ui="{ footer: 'justify-end', content: 'max-w-2xl' }">
    <template #body>
      <div
        v-if="reversosAusentes.length === 0"
        class="flex flex-col items-center gap-2 py-10 text-center">
        <UIcon name="i-lucide-circle-check" class="size-7 text-success" />
        <p class="text-sm font-medium text-toned">Nenhum reverso ausente</p>
        <p class="max-w-xs text-xs text-dimmed">
          Todos os registros A e AAAA com zona reversa cadastrada já têm PTR.
        </p>
      </div>

      <div v-else class="space-y-3">
        <div class="flex items-center justify-between gap-3">
          <UCheckbox
            :model-value="selecaoReversos"
            :label="`${reversosSelecionados.length} de ${reversosAusentes.length} selecionado(s)`"
            :ui="{ label: 'text-sm text-toned' }"
            @update:model-value="toggleTodosReversos" />
        </div>

        <div
          class="max-h-80 divide-y divide-default overflow-y-auto rounded-lg border border-default">
          <label
            v-for="reverso in reversosAusentes"
            :key="reverso.nomeReverso"
            class="flex cursor-pointer items-center justify-between gap-3 p-2.5 transition-colors hover:bg-elevated/60">
            <div class="flex min-w-0 items-center gap-2.5">
              <UCheckbox
                :model-value="reversosSelecionados.includes(reverso.nomeReverso)"
                :aria-label="`Selecionar o reverso de ${reverso.ip}`"
                @update:model-value="toggleReverso(reverso.nomeReverso, $event === true)" />

              <div class="min-w-0">
                <p :title="reverso.name" class="truncate data text-sm font-medium text-highlighted">
                  {{ reverso.name }}
                </p>
                <p :title="reverso.nomeReverso" class="truncate data text-xs text-dimmed">
                  {{ reverso.nomeReverso }}
                </p>
              </div>
            </div>

            <div class="flex shrink-0 items-center gap-2">
              <span
                class="inline-flex items-center rounded-md px-2 py-0.5 data text-[11px] font-semibold ring-1 ring-inset"
                :class="typeMeta(reverso.type).chip">
                {{ reverso.type }}
              </span>
              <span class="data text-xs text-muted tnum">{{ reverso.ip }}</span>
            </div>
          </label>
        </div>
      </div>
    </template>

    <template #footer>
      <UButton
        label="Fechar"
        color="neutral"
        variant="ghost"
        :loading="isLoading"
        @click="
          () => {
            modalReversos = false
          }
        " />
      <UButton
        v-if="reversosAusentes.length > 0"
        :label="
          reversosSelecionados.length === 1
            ? 'Criar 1 reverso'
            : `Criar ${reversosSelecionados.length} reversos`
        "
        icon="i-lucide-check"
        :loading="isLoading"
        :disabled="reversosSelecionados.length === 0"
        @click="createReversos" />
    </template>
  </UModal>

  <UModal
    v-model:open="modalOrfaos"
    title="Reversos órfãos"
    :description="`Registros PTR de ${zoneLabel} sem nenhum A ou AAAA apontando para o IP`"
    :ui="{ footer: 'justify-end', content: 'max-w-2xl' }">
    <template #body>
      <div
        v-if="reversosOrfaos.length === 0"
        class="flex flex-col items-center gap-2 py-10 text-center">
        <UIcon name="i-lucide-circle-check" class="size-7 text-success" />
        <p class="text-sm font-medium text-toned">Nenhum reverso órfão</p>
        <p class="max-w-xs text-xs text-dimmed">
          Todos os PTR desta zona têm um registro direto correspondente.
        </p>
      </div>

      <div v-else class="space-y-3">
        <p class="text-sm text-toned">
          <span class="data text-highlighted tnum">{{ reversosOrfaos.length }}</span>
          {{ reversosOrfaos.length === 1 ? 'PTR sem registro direto' : 'PTRs sem registro direto' }}
        </p>

        <div
          class="max-h-80 divide-y divide-default overflow-y-auto rounded-lg border border-default">
          <div
            v-for="orfao in reversosOrfaos"
            :key="orfao.nomeReverso"
            class="flex items-center justify-between gap-3 p-2.5">
            <div class="min-w-0">
              <p :title="orfao.alvo" class="truncate data text-sm font-medium text-highlighted">
                {{ orfao.alvo }}
              </p>
              <p :title="orfao.nomeReverso" class="truncate data text-xs text-dimmed">
                {{ orfao.nomeReverso }}
              </p>
            </div>

            <div class="flex shrink-0 items-center gap-2">
              <span class="data text-xs text-muted tnum">{{ orfao.ip }}</span>

              <UPopover>
                <UButton
                  icon="i-lucide-trash-2"
                  color="error"
                  variant="ghost"
                  size="sm"
                  :loading="isLoading"
                  :aria-label="`Excluir o PTR de ${orfao.ip}`" />

                <template #content="{ close }">
                  <div class="w-64 space-y-3 p-3">
                    <p class="text-sm text-toned">
                      Excluir o PTR de
                      <span class="data text-highlighted">{{ orfao.ip }}</span>
                      ?
                    </p>
                    <div class="flex justify-end gap-2">
                      <UButton
                        color="neutral"
                        variant="ghost"
                        size="sm"
                        label="Cancelar"
                        @click="close()" />
                      <UButton
                        color="error"
                        size="sm"
                        label="Excluir"
                        :loading="isLoading"
                        @click="deleteOrfao(orfao.nomeReverso, close)" />
                    </div>
                  </div>
                </template>
              </UPopover>
            </div>
          </div>
        </div>
      </div>
    </template>

    <template #footer>
      <UButton
        label="Fechar"
        color="neutral"
        variant="ghost"
        :loading="isLoading"
        @click="
          () => {
            modalOrfaos = false
          }
        " />
    </template>
  </UModal>

  <Users v-if="isAdmin" v-model:open="modalUsers" :zone="zoneId" />
</template>
