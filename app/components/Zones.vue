<script setup lang="ts">
  import type { BadgeProps, DropdownMenuItem, TableColumn, TableRow } from '@nuxt/ui'
  import type { Row } from '@tanstack/vue-table'
  import { getPaginationRowModel } from '@tanstack/vue-table'
  import { safeParse } from 'valibot'

  import { UBadge, UButton, UDropdownMenu, UIcon } from '#components'

  const { user } = useUser()
  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()

  const zoneId = defineModel<string>('zoneId', { required: true })
  const type = defineModel<string>('type', { required: true })
  const nivel = defineModel<string>('nivel', { required: true })

  const { data, refresh } = await useFetch<ZonesResponse>('/server/api/zones', {
    method: 'GET',
    query: { type },
  })

  const globalFilter = ref('')

  const table = useTemplateRef('table')

  const pagination = ref({ pageIndex: 0, pageSize: 10 })

  watch(globalFilter, () => {
    pagination.value.pageIndex = 0
  })

  const columns: TableColumn<ZoneFetch>[] = [
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
    },
    {
      accessorKey: 'serial',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Serial',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
    },
    {
      accessorKey: 'nivel',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Nivel',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
      cell: ({ row }) => {
        const { nivel: level } = row.original
        let color: BadgeProps['color'] = 'success'

        if (level === 'ADMINISTRADOR') {
          color = 'error'
        } else if (level === 'LEITURA') {
          color = 'warning'
        }
        return h(UBadge, { color, variant: 'soft', class: 'uppercase' }, () => level)
      },
    },
    {
      accessorKey: 'dnssec',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'DNSSEC',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
      cell: ({ row }) =>
        row.original.dnssec
          ? // oxlint-disable-next-line antfu/consistent-list-newline
            h(UBadge, { variant: 'soft', color: 'success' }, () =>
              h('div', { class: 'flex items-center space-x-2' }, [
                h(UIcon, { name: 'i-lucide-shield-check', class: 'text-green-500' }),
                h('span', 'Ativo'),
              ]),
            )
          : // oxlint-disable-next-line antfu/consistent-list-newline
            h(UBadge, { variant: 'soft', color: 'error' }, () =>
              h('div', { class: 'flex items-center space-x-2' }, [
                h(UIcon, { name: 'i-lucide-shield-off', class: 'text-red-500' }),
                h('span', 'Desativado'),
              ]),
            ),
    },
    {
      id: 'actions',
      cell: ({ row }) =>
        h(
          'div',
          { class: 'text-right' },
          h(
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
          ),
        ),
    },
  ]

  function onSelect(e: Event, row: TableRow<ZoneFetch>): void {
    zoneId.value = row.original.name
    nivel.value = row.original.nivel
  }

  const modalDelete = ref(false)
  const idDelete = ref('')
  const confirmIdDelete = ref('')

  const { copy } = useClipboard()

  function getRowItems(row: Row<ZoneFetch>): DropdownMenuItem[] {
    return [
      { type: 'label', label: `${row.original.name} actions` },
      {
        label: 'Copy zone',
        icon: 'i-lucide-copy',
        onSelect(): void {
          copy(row.original.name)
          toast.add({
            title: 'Zone copied to clipboard!',
            color: 'success',
            icon: 'i-lucide-circle-check',
          })
        },
      },
      { type: 'separator' },
      {
        label: 'View zone',
        icon: 'i-lucide-eye',
        onSelect(): void {
          zoneId.value = row.original.name
          nivel.value = row.original.nivel
        },
      },
      ...(user.value.level === 'admin'
        ? ([
            {
              label: 'Delete zone',
              icon: 'i-lucide-trash',
              color: 'error',
              onSelect(): void {
                idDelete.value = row.original.name
                modalDelete.value = true
              },
            },
          ] satisfies DropdownMenuItem[])
        : []),
    ]
  }

  watch(modalDelete, (nv) => {
    if (!nv) {
      idDelete.value = ''
      confirmIdDelete.value = ''
    }
  })

  async function deleteZone(): Promise<void> {
    start()

    if (confirmIdDelete.value !== idDelete.value) {
      toast.add({
        title: 'The zone ID entered does not match.',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
      return finish({ error: true })
    }

    if (idDelete.value === '') {
      toast.add({ title: 'Zone ID is required.', icon: 'i-lucide-shield-alert', color: 'error' })
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/zone', {
      method: 'DELETE',
      query: { id: idDelete.value },
    }).catch((error) => {
      toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
    })

    if (!res) return finish({ error: true })

    if (!res.message) {
      toast.add({
        title: 'An unknown error occurred',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
      return finish({ error: true })
    }

    finish({ force: true })
    refresh()
    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    modalDelete.value = false
  }

  const modal = ref(false)

  const state = ref<ZoneSchemaType>({
    domain: '',
    type: '',
    soa: {
      startOfAuthority: '',
      email: '',
      refresh: 3600,
      retry: 600,
      expire: 604_800,
      negativeCacheTtl: 86_400,
    },
  })

  watch(modal, (nv) => {
    if (!nv) {
      state.value = {
        domain: '',
        type: '',
        soa: {
          startOfAuthority: '',
          email: '',
          refresh: 3600,
          retry: 600,
          expire: 604_800,
          negativeCacheTtl: 86_400,
        },
      }
    }
  })

  async function createZone(): Promise<void> {
    start()

    state.value.type = type.value

    const body = safeParse(ZoneSchema, state.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/zone', { method: 'PUT', body: body.output }).catch(
      (error) => {
        toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
      },
    )

    if (!res) return finish({ error: true })

    if (!res.message) {
      toast.add({
        title: 'An unknown error occurred',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
      return finish({ error: true })
    }

    finish({ force: true })
    refresh()
    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    modal.value = false
  }
</script>

<template>
  <div class="my-6 flex items-center justify-between">
    <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">Zonas</h1>
    <div class="grid grid-cols-1 gap-2 md:flex md:items-center md:gap-4">
      <UButton
        v-if="user.level === 'admin'"
        color="info"
        :label="
          type === 'normal'
            ? 'Criar Zona'
            : type === 'reverse'
              ? 'Criar Zona Reversa'
              : 'Criar Zona Reversa IPv6'
        "
        icon="i-lucide-plus"
        variant="soft"
        @click="modal = true" />
      <UInput v-model="globalFilter" placeholder="Pesquisar zonas..." />
    </div>
  </div>

  <ClientOnly>
    <UTable
      ref="table"
      v-model:global-filter="globalFilter"
      v-model:pagination="pagination"
      :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
      :data="data?.zones"
      :columns="columns"
      @select="onSelect" />
  </ClientOnly>

  <div
    v-if="data?.zones && data.zones.length > pagination.pageSize"
    class="flex justify-center border-t border-default pt-4">
    <UPagination
      active-color="info"
      color="info"
      active-variant="subtle"
      :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1"
      :items-per-page="table?.tableApi?.getState().pagination.pageSize"
      :total="table?.tableApi?.getFilteredRowModel().rows.length"
      @update:page="(p) => table?.tableApi?.setPageIndex(p - 1)" />
  </div>

  <UModal
    v-model:open="modal"
    title="Criar Zona"
    description="Criar uma nova zona no Servidor DNS"
    :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :schema="ZoneSchema" :state="state" class="space-y-4">
        <UFormField label="Domínio" name="domain">
          <UInput
            v-model="state.domain"
            icon="i-lucide-computer"
            class="w-full"
            placeholder="Ex: example.com" />
        </UFormField>

        <USeparator label="Configurações do Registro Start of Authority (SOA)" />

        <UFormField label="Start of Authority" name="soa.startOfAuthority">
          <UInput
            v-model="state.soa.startOfAuthority"
            icon="i-lucide-shield-check"
            class="w-full"
            placeholder="Ex: ns1.example.com" />
        </UFormField>
        <UFormField label="Email" name="soa.email">
          <UInput
            v-model="state.soa.email"
            icon="i-lucide-mail"
            class="w-full"
            placeholder="Ex: hostmaster.example.com" />
        </UFormField>
        <UFormField label="Refresh" name="soa.refresh">
          <UInputNumber
            v-model="state.soa.refresh"
            :min="0"
            icon="i-lucide-refresh-cw"
            class="w-full"
            placeholder="3600" />
        </UFormField>
        <UFormField label="Retry" name="soa.retry">
          <UInputNumber
            v-model="state.soa.retry"
            :min="0"
            icon="i-lucide-clock"
            class="w-full"
            placeholder="600" />
        </UFormField>
        <UFormField label="Expire" name="soa.expire">
          <UInputNumber
            v-model="state.soa.expire"
            :min="0"
            icon="i-lucide-hourglass"
            class="w-full"
            placeholder="604800" />
        </UFormField>
        <UFormField label="Negative Cache TTL" name="soa.negativeCacheTtl">
          <UInputNumber
            v-model="state.soa.negativeCacheTtl"
            :min="0"
            icon="i-lucide-timer"
            class="w-full"
            placeholder="3600" />
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <UButton label="Cancelar" :loading="isLoading" variant="outline" @click="modal = false" />
      <UButton label="Confirmar" :loading="isLoading" @click="createZone" />
    </template>
  </UModal>

  <UModal
    v-model:open="modalDelete"
    title="Aviso"
    description="Você está prestes a deletar uma zona, esta ação não pode ser desfeita."
    :ui="{ footer: 'justify-end' }">
    <template #body>
      <p class="dark:text-gray-200">
        Se você tem certeza que deseja continuar, escreva o ID da zona abaixo
        <span class="font-bold">'{{ idDelete }}'</span>.
      </p>
      <UInput
        v-model="confirmIdDelete"
        class="mt-2 w-full"
        color="error"
        placeholder="ID da Zona" />
    </template>

    <template #footer>
      <UButton
        label="Cancelar"
        :loading="isLoading"
        variant="outline"
        @click="modalDelete = false" />
      <UButton
        label="Confirmar"
        color="error"
        :loading="isLoading"
        :disabled="confirmIdDelete !== idDelete"
        @click="deleteZone" />
    </template>
  </UModal>
</template>
