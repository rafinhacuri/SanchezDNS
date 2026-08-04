<script setup lang="ts">
  import type { DropdownMenuItem, TableColumn, TableRow } from '@nuxt/ui'
  import type { Column, Row } from '@tanstack/vue-table'
  import { getPaginationRowModel } from '@tanstack/vue-table'
  import { safeParse } from 'valibot'
  import type { VNode } from 'vue'

  import { UButton, UDropdownMenu } from '#components'

  const { user } = useUser()
  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()
  const { zoneSchema } = useZonesSchema()

  const props = defineProps({
    type: { type: String, required: true },
  })

  const type = computed(() => props.type)

  const { data, refresh } = await useApi<ZonesResponse>('/zones', {
    method: 'GET',
    query: { type },
  })

  const globalFilter = ref('')

  const table = useTemplateRef('table')

  const pagination = ref({ pageIndex: 0, pageSize: 10 })

  watch(globalFilter, () => {
    pagination.value.pageIndex = 0
  })

  function sortableHeader(column: Column<ZoneFetch>, label: string): VNode {
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

  const columns: TableColumn<ZoneFetch>[] = [
    {
      accessorKey: 'name',
      header: ({ column }) => sortableHeader(column, 'Nome'),
      cell: ({ row }) =>
        h(
          'span',
          { class: 'data font-medium text-highlighted', title: row.original.name },
          row.original.name,
        ),
    },
    {
      accessorKey: 'serial',
      header: ({ column }) => sortableHeader(column, 'Serial'),
      cell: ({ row }) => h('span', { class: 'data text-sm text-muted' }, row.original.serial),
    },
    {
      accessorKey: 'nivel',
      header: ({ column }) => sortableHeader(column, 'Nível'),
      cell: ({ row }) => {
        const { nivel: level } = row.original
        let chip = 'bg-emerald-500/10 text-emerald-600 ring-emerald-500/25 dark:text-emerald-400'
        if (level === 'ADMINISTRADOR') {
          chip = 'bg-rose-500/10 text-rose-600 ring-rose-500/25 dark:text-rose-400'
        } else if (level === 'LEITURA') {
          chip = 'bg-amber-500/10 text-amber-600 ring-amber-500/25 dark:text-amber-400'
        }

        return h(
          'span',
          {
            class: `inline-flex items-center rounded-md px-2 py-0.5 text-[11px] font-semibold uppercase tracking-wide ring-1 ring-inset ${chip}`,
          },
          level,
        )
      },
    },
    {
      accessorKey: 'dnssec',
      header: ({ column }) => sortableHeader(column, 'DNSSEC'),
      cell: ({ row }) => {
        const active = row.original.dnssec

        return h('div', { class: 'flex items-center gap-1.5' }, [
          h('span', {
            class: `size-1.5 rounded-full ${active ? 'bg-emerald-500' : 'bg-slate-400 dark:bg-slate-600'}`,
          }),
          h(
            'span',
            { class: active ? 'text-sm text-toned' : 'text-sm text-dimmed' },
            active ? 'Ativo' : 'Desativado',
          ),
        ])
      },
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
              'aria-label': 'Ações da zona',
            },
            () =>
              h(UButton, {
                icon: 'i-lucide-ellipsis-vertical',
                color: 'neutral',
                variant: 'ghost',
                size: 'sm',
                class: 'ml-auto',
                'aria-label': `Ações da zona ${row.original.name}`,
              }),
          ),
        ),
    },
  ]

  async function openZone(zone: ZoneFetch): Promise<void> {
    await navigateTo(`/zonas/${zone.name}`)
  }

  function onSelect(e: Event, row: TableRow<ZoneFetch>): void {
    void openZone(row.original)
  }

  const modalDelete = ref(false)
  const idDelete = ref('')
  const confirmIdDelete = ref('')

  const { copy } = useClipboard()

  function getRowItems(row: Row<ZoneFetch>): DropdownMenuItem[] {
    return [
      { type: 'label', label: row.original.name },
      {
        label: 'Copiar zona',
        icon: 'i-lucide-copy',
        onSelect(): void {
          copy(row.original.name)
          toast.add({
            title: 'Zona copiada para a área de transferência',
            color: 'success',
            icon: 'i-lucide-circle-check',
          })
        },
      },
      { type: 'separator' },
      {
        label: 'Abrir zona',
        icon: 'i-lucide-arrow-right',
        onSelect(): void {
          openZone(row.original)
        },
      },
      ...(user.value.level === 'admin'
        ? [
            { type: 'separator' as const },
            {
              label: 'Deletar zona',
              icon: 'i-lucide-trash-2',
              color: 'error' as const,
              onSelect(): void {
                idDelete.value = row.original.name
                modalDelete.value = true
              },
            },
          ]
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
        title: 'O ID da zona informado não confere.',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
      return finish({ error: true })
    }

    if (idDelete.value === '') {
      toast.add({
        title: 'O ID da zona é obrigatório.',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
      return finish({ error: true })
    }

    const res = await $api('/zone', {
      method: 'DELETE',
      query: { id: idDelete.value },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao deletar zona',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    finish({ force: true })
    await refresh()
    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    modalDelete.value = false
  }

  const modal = ref(false)

  const state = ref<ZoneSchema>({
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

    const body = safeParse(zoneSchema, state.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/zone', { method: 'PUT', body: body.output }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao criar zona',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    finish({ force: true })
    await refresh()
    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    modal.value = false
  }

  const createLabel = computed(() => {
    if (type.value === 'reverse') return 'Criar zona reversa'
    if (type.value === 'reverse-ipv6') return 'Criar zona reversa IPv6'
    return 'Criar zona'
  })

  const totalZones = computed(() => data.value?.zones?.length ?? 0)
</script>

<template>
  <div class="overflow-hidden rounded-xl border border-default bg-default">
    <div
      class="flex flex-col gap-3 border-b border-default bg-muted/40 p-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-2.5">
        <UIcon name="i-lucide-database" class="size-4 shrink-0 text-dimmed" />
        <p class="text-sm font-medium text-toned">
          <span class="data text-highlighted tnum">{{ totalZones }}</span>
          {{ totalZones === 1 ? 'zona' : 'zonas' }}
        </p>
      </div>

      <div class="flex items-center gap-2">
        <UInput
          v-model="globalFilter"
          icon="i-lucide-search"
          placeholder="Filtrar zonas..."
          class="w-full sm:w-64"
          :ui="{ base: 'data' }" />

        <UButton
          v-if="user.level === 'admin'"
          :label="createLabel"
          icon="i-lucide-plus"
          class="shrink-0"
          @click="
            () => {
              modal = true
            }
          " />
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
        :ui="{
          tr: 'cursor-pointer transition-colors duration-150 hover:bg-elevated/60',
          th: 'py-2',
          td: 'py-2.5',
        }"
        @select="onSelect">
        <template #empty>
          <div class="flex flex-col items-center gap-2 py-12 text-center">
            <UIcon name="i-lucide-database-zap" class="size-7 text-dimmed" />
            <p class="text-sm font-medium text-toned">Nenhuma zona encontrada</p>
            <p class="max-w-xs text-xs text-dimmed">
              {{
                globalFilter
                  ? 'Nenhum resultado para esse filtro. Tente outro termo.'
                  : 'Ainda não há zonas cadastradas para este tipo.'
              }}
            </p>
          </div>
        </template>
      </UTable>
    </ClientOnly>

    <div
      v-if="data?.zones && data.zones.length > pagination.pageSize"
      class="flex justify-center border-t border-default bg-muted/40 p-3">
      <UPagination
        active-variant="subtle"
        :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1"
        :items-per-page="table?.tableApi?.getState().pagination.pageSize"
        :total="table?.tableApi?.getFilteredRowModel().rows.length"
        @update:page="(p) => table?.tableApi?.setPageIndex(p - 1)" />
    </div>
  </div>

  <UModal
    v-model:open="modal"
    :title="createLabel"
    description="Uma nova zona precisa de um domínio e das configurações do registro SOA."
    :ui="{ content: 'max-w-2xl', footer: 'justify-end' }">
    <template #body>
      <UForm :schema="zoneSchema" :state="state" class="space-y-6">
        <UFormField label="Domínio" name="domain" required>
          <UInput
            v-model="state.domain"
            icon="i-lucide-globe"
            class="w-full"
            placeholder="example.com"
            :ui="{ base: 'data' }" />
        </UFormField>

        <div class="space-y-4 rounded-lg border border-default bg-muted/40 p-4">
          <div class="flex items-center gap-2">
            <UIcon name="i-lucide-file-cog" class="size-4 text-dimmed" />
            <h3 class="text-sm font-semibold text-highlighted">Start of Authority (SOA)</h3>
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <UFormField label="Servidor primário" name="soa.startOfAuthority" required>
              <UInput
                v-model="state.soa.startOfAuthority"
                icon="i-lucide-server"
                class="w-full"
                placeholder="ns1.example.com"
                :ui="{ base: 'data' }" />
            </UFormField>

            <UFormField label="Email do responsável" name="soa.email" required>
              <UInput
                v-model="state.soa.email"
                icon="i-lucide-mail"
                class="w-full"
                placeholder="hostmaster.example.com"
                :ui="{ base: 'data' }" />
            </UFormField>

            <UFormField label="Refresh" name="soa.refresh" hint="segundos">
              <UInputNumber
                v-model="state.soa.refresh"
                :min="0"
                class="w-full"
                placeholder="3600" />
            </UFormField>

            <UFormField label="Retry" name="soa.retry" hint="segundos">
              <UInputNumber v-model="state.soa.retry" :min="0" class="w-full" placeholder="600" />
            </UFormField>

            <UFormField label="Expire" name="soa.expire" hint="segundos">
              <UInputNumber
                v-model="state.soa.expire"
                :min="0"
                class="w-full"
                placeholder="604800" />
            </UFormField>

            <UFormField label="Negative Cache TTL" name="soa.negativeCacheTtl" hint="segundos">
              <UInputNumber
                v-model="state.soa.negativeCacheTtl"
                :min="0"
                class="w-full"
                placeholder="86400" />
            </UFormField>
          </div>
        </div>
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
            modal = false
          }
        " />
      <UButton label="Criar zona" icon="i-lucide-check" :loading="isLoading" @click="createZone" />
    </template>
  </UModal>

  <UModal
    v-model:open="modalDelete"
    title="Deletar zona"
    description="Esta ação é permanente e não pode ser desfeita."
    :ui="{ footer: 'justify-end' }">
    <template #body>
      <div class="space-y-4">
        <div
          class="flex gap-3 rounded-lg border border-rose-500/20 bg-rose-500/5 p-3 text-sm text-toned">
          <UIcon name="i-lucide-triangle-alert" class="mt-0.5 size-4 shrink-0 text-rose-500" />
          <p>Todos os registros desta zona serão removidos junto com ela.</p>
        </div>

        <UFormField label="Confirme digitando o nome da zona">
          <template #hint>
            <span class="data text-xs text-highlighted">{{ idDelete }}</span>
          </template>
          <UInput
            v-model="confirmIdDelete"
            class="w-full"
            color="error"
            autocomplete="off"
            :placeholder="idDelete"
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
        label="Deletar zona"
        icon="i-lucide-trash-2"
        color="error"
        :loading="isLoading"
        :disabled="confirmIdDelete !== idDelete"
        @click="deleteZone" />
    </template>
  </UModal>
</template>
