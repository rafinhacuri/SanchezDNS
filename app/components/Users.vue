<script setup lang="ts">
  import type { TableColumn } from '@nuxt/ui'
  import type { Column } from '@tanstack/vue-table'
  import { getPaginationRowModel } from '@tanstack/vue-table'
  import { safeParse } from 'valibot'
  import type { VNode } from 'vue'

  import { UButton, UPopover } from '#components'

  const baseUrl = useApiUrl()

  const tableUsers = useTemplateRef('tableUsers')

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()
  const { insertUserSchema } = useUserSchema()

  const props = defineProps({
    zone: { type: String, required: true },
  })

  const zoneId = computed(() => props.zone)

  const modalUsers = defineModel<boolean>('open', { default: false })

  const { data: userData, refresh: refreshUsers } = await useApi<User[]>('/users', {
    method: 'GET',
    query: { zona: zoneId },
  })

  const { data: members } = await useApi<{ email: string; nome: string }[]>('/members', {
    method: 'GET',
  })

  const roleOptions = ['escrita', 'leitura']

  const stateUser = ref<InsertUserSchema>({ email: '', permissao: 'leitura', id: '', zona: '' })
  const isEditingUser = ref(false)

  async function addUser(): Promise<void> {
    start()

    stateUser.value.zona = zoneId.value

    const body = safeParse(insertUserSchema, stateUser.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/user', {
      method: body.output.id ? 'PATCH' : 'POST',
      body: body.output,
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao adicionar usuário',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    isEditingUser.value = false
    await refreshUsers()
    stateUser.value = { email: '', permissao: 'leitura', id: '', zona: '' }
    finish()
  }

  watch(modalUsers, (nv) => {
    if (!nv) stateUser.value = { email: '', permissao: 'leitura', id: '', zona: '' }
  })

  async function deleteUser(zona: string, id: string): Promise<void> {
    start()

    const res = await $api('/user', {
      method: 'DELETE',
      body: { zona, id },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao deletar usuário',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    await refreshUsers()
    finish()
  }

  const paginationUsers = ref({ pageIndex: 0, pageSize: 5 })
  const globalFilterUsers = ref('')

  watch(globalFilterUsers, () => {
    paginationUsers.value.pageIndex = 0
  })

  function sortableHeader(column: Column<User>, label: string): VNode {
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

  const columnsUsers: TableColumn<User>[] = [
    {
      accessorKey: 'email',
      header: ({ column }) => sortableHeader(column, 'Usuário'),
      cell: ({ row }) =>
        h('div', { class: 'flex items-center gap-2.5' }, [
          h('img', {
            src: `${baseUrl}/file/${row.original.email}`,
            alt: row.original.email,
            class: 'size-7 rounded-full object-cover ring-1 ring-default',
          }),
          h('span', { class: 'data text-sm font-medium text-highlighted' }, row.original.email),
        ]),
    },
    {
      accessorKey: 'permissao',
      header: ({ column }) => sortableHeader(column, 'Permissão'),
      cell: ({ row }) => {
        const write = row.original.permissao === 'escrita'

        return h(
          'span',
          {
            class: `inline-flex items-center gap-1.5 rounded-md px-2 py-0.5 text-[11px] font-semibold uppercase tracking-wide ring-1 ring-inset ${
              write
                ? 'bg-blue-500/10 text-blue-600 ring-blue-500/25 dark:text-blue-400'
                : 'bg-slate-500/10 text-slate-600 ring-slate-500/25 dark:text-slate-300'
            }`,
          },
          write ? 'Escrita' : 'Leitura',
        )
      },
    },
    {
      accessorKey: 'actions',
      header: (): VNode => h('span', { class: 'sr-only' }, 'Ações'),
      cell: ({ row }) =>
        h('div', { class: 'flex items-center justify-end gap-1.5' }, [
          !stateUser.value.id &&
            stateUser.value.id !== row.original.id &&
            h(UButton, {
              icon: 'i-lucide-pencil',
              color: 'neutral',
              variant: 'ghost',
              size: 'sm',
              'aria-label': `Editar ${row.original.email}`,
              onClick: () => {
                isEditingUser.value = true
                stateUser.value = {
                  id: row.original.id,
                  email: row.original.email,
                  zona: row.original.zona,
                  permissao: row.original.permissao === 'escrita' ? 'escrita' : 'leitura',
                }
              },
            }),

          stateUser.value.id === row.original.id &&
            h(UButton, {
              icon: 'i-lucide-x',
              color: 'neutral',
              variant: 'ghost',
              size: 'sm',
              'aria-label': 'Cancelar edição',
              onClick: () => {
                isEditingUser.value = false
                stateUser.value = { email: '', permissao: 'leitura', id: '', zona: '' }
              },
            }),

          !stateUser.value.id &&
            h(UPopover, null, {
              default: () =>
                h(UButton, {
                  icon: 'i-lucide-user-minus',
                  color: 'error',
                  variant: 'ghost',
                  size: 'sm',
                  'aria-label': `Remover ${row.original.email}`,
                }),
              content: ({ close }: { close: () => void }) =>
                h('div', { class: 'w-64 space-y-3 p-3' }, [
                  h(
                    'p',
                    { class: 'text-sm text-toned' },
                    `Remover ${row.original.email} desta zona?`,
                  ),
                  h('div', { class: 'flex justify-end gap-2' }, [
                    h(UButton, {
                      color: 'neutral',
                      variant: 'ghost',
                      size: 'sm',
                      label: 'Cancelar',
                      onClick: () => close(),
                    }),
                    h(UButton, {
                      color: 'error',
                      size: 'sm',
                      label: 'Remover',
                      onClick: () => {
                        deleteUser(row.original.zona, row.original.id)
                        close()
                      },
                    }),
                  ]),
                ]),
            }),
        ]),
    },
  ]

  const totalUsers = computed(() => userData.value?.length ?? 0)
</script>

<template>
  <UModal
    v-model:open="modalUsers"
    title="Usuários da zona"
    :description="`Quem pode ler e escrever registros em ${zoneId}`"
    :ui="{ footer: 'justify-end', content: 'max-w-3xl' }">
    <template #body>
      <div class="space-y-5">
        <div class="space-y-3 rounded-lg border border-default bg-muted/40 p-4">
          <p class="text-sm font-semibold text-highlighted">
            {{ isEditingUser ? 'Editar permissão' : 'Adicionar usuário' }}
          </p>

          <div class="flex flex-col gap-3 sm:flex-row">
            <USelectMenu
              v-model="stateUser.email"
              :disabled="isEditingUser"
              label-key="nome"
              value-key="email"
              :items="members || []"
              icon="i-lucide-search"
              placeholder="Selecione o usuário..."
              class="w-full sm:flex-1"
              :ui="{ base: 'data' }">
              <template #item-label="{ item }">
                <div class="flex items-center gap-2.5">
                  <img
                    :src="`${baseUrl}/file/${item.email}`"
                    :alt="item.nome"
                    class="size-6 rounded-full object-cover ring-1 ring-default" />
                  <div class="min-w-0">
                    <p class="truncate text-sm font-medium text-highlighted">{{ item.nome }}</p>
                    <p class="truncate data text-xs text-dimmed">{{ item.email }}</p>
                  </div>
                </div>
              </template>
            </USelectMenu>

            <USelect
              v-model="stateUser.permissao"
              :items="roleOptions"
              icon="i-lucide-shield-check"
              placeholder="Permissão..."
              class="w-full sm:w-44" />

            <UButton
              :label="stateUser.id ? 'Salvar' : 'Adicionar'"
              :icon="stateUser.id ? 'i-lucide-check' : 'i-lucide-user-plus'"
              class="justify-center sm:w-auto"
              :loading="isLoading"
              @click="addUser" />
          </div>
        </div>

        <div class="overflow-hidden rounded-lg border border-default">
          <div
            class="flex items-center justify-between gap-3 border-b border-default bg-muted/40 p-2.5">
            <p class="text-xs font-medium text-muted">
              <span class="data text-highlighted tnum">{{ totalUsers }}</span>
              {{ totalUsers === 1 ? 'usuário' : 'usuários' }}
            </p>
            <UInput
              v-model="globalFilterUsers"
              icon="i-lucide-search"
              size="sm"
              placeholder="Filtrar..."
              class="w-48"
              :ui="{ base: 'data' }" />
          </div>

          <ClientOnly>
            <UTable
              ref="tableUsers"
              v-model:global-filter="globalFilterUsers"
              v-model:pagination="paginationUsers"
              :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
              :data="userData"
              :columns="columnsUsers"
              :ui="{
                tr: 'transition-colors duration-150 hover:bg-elevated/60',
                th: 'py-2',
                td: 'py-2',
              }">
              <template #empty>
                <div class="flex flex-col items-center gap-1.5 py-10 text-center">
                  <UIcon name="i-lucide-users" class="size-6 text-dimmed" />
                  <p class="text-sm font-medium text-toned">Nenhum usuário nesta zona</p>
                  <p class="text-xs text-dimmed">Adicione alguém no campo acima.</p>
                </div>
              </template>
            </UTable>
          </ClientOnly>

          <div
            v-if="userData && userData.length > paginationUsers.pageSize"
            class="flex justify-center border-t border-default bg-muted/40 p-2.5">
            <UPagination
              size="sm"
              active-variant="subtle"
              :default-page="(tableUsers?.tableApi?.getState().pagination.pageIndex || 0) + 1"
              :items-per-page="tableUsers?.tableApi?.getState().pagination.pageSize"
              :total="tableUsers?.tableApi?.getFilteredRowModel().rows.length"
              @update:page="(p) => tableUsers?.tableApi?.setPageIndex(p - 1)" />
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
            modalUsers = false
          }
        " />
    </template>
  </UModal>
</template>
