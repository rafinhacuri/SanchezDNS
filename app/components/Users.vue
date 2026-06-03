<script setup lang="ts">
  import type { TableColumn } from '@nuxt/ui'
  import { getPaginationRowModel } from '@tanstack/vue-table'
  import { safeParse } from 'valibot'

  import { UBadge, UButton, UPopover } from '#components'

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()

  const zoneId = defineModel<string>('zoneId', { required: true })
  const modalUsers = defineModel<boolean>('open', { default: false })

  const { data: userData, refresh: refreshUsers } = await useFetch<User[]>('/server/api/users', {
    method: 'GET',
    query: { zona: zoneId },
  })

  const { data: members } = await useFetch<{ email: string; nome: string }[]>(
    '/server/api/members',
    {
      method: 'GET',
    },
  )

  const roleOptions = ['escrita', 'leitura']

  const stateUser = ref<InsertUserType>({ email: '', permissao: 'leitura', id: '', zona: '' })
  const isEditingUser = ref(false)

  async function addUser(): Promise<void> {
    start()

    stateUser.value.zona = zoneId.value

    const body = safeParse(InsertUserSchema, stateUser.value)

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/user', {
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

    const res = await $fetch<GoRes>('/server/api/user', {
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

  const tableUsers = useTemplateRef('tableUsers')

  const paginationUsers = ref({ pageIndex: 0, pageSize: 5 })
  const globalFilterUsers = ref('')

  watch(globalFilterUsers, () => {
    paginationUsers.value.pageIndex = 0
  })

  const columnsUsers: TableColumn<User>[] = [
    {
      accessorKey: 'email',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Usuário',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
      cell: ({ row }) =>
        h('div', { class: 'flex items-center gap-3' }, [
          h('img', {
            src: `/server/api/file/${row.original.email}`,
            alt: row.original.email,
            class: 'size-6 cursor-pointer rounded-full',
          }),
          h('p', {}, row.original.email),
        ]),
    },
    {
      accessorKey: 'permissao',
      header: ({ column }) => {
        const isSorted = column.getIsSorted()
        let icon = 'i-heroicons-arrows-up-down'
        if (isSorted === 'asc') icon = 'i-heroicons-bars-arrow-up'
        else if (isSorted === 'desc') icon = 'i-heroicons-bars-arrow-down'
        return h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          label: 'Permissão',
          icon,
          class: '-mx-2.5',
          onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
        })
      },
      cell: ({ row }) =>
        row.original.permissao === 'escrita'
          ? h(UBadge, { variant: 'outline', label: 'Escrita' })
          : h(UBadge, { color: 'neutral', variant: 'outline', label: 'Leitura' }),
    },
    {
      accessorKey: 'actions',
      header: 'Ações',
      cell: ({ row }) =>
        h('div', { class: 'space-x-2' }, [
          !stateUser.value.id &&
            stateUser.value.id !== row.original.id &&
            h(UButton, {
              icon: 'i-lucide-user-pen',
              variant: 'outline',
              onClick: () => {
                isEditingUser.value = true
                stateUser.value = { ...row.original }
              },
            }),

          stateUser.value.id === row.original.id &&
            h(UButton, {
              icon: 'i-lucide-x',
              variant: 'outline',
              onClick: () => {
                isEditingUser.value = false
                stateUser.value = { email: '', permissao: 'leitura', id: '', zona: '' }
              },
            }),

          !stateUser.value.id &&
            h(UPopover, null, {
              default: () =>
                h(UButton, { icon: 'i-lucide-user-minus', color: 'error', variant: 'outline' }),
              content: ({ close }: { close: () => void }) =>
                h('div', { class: 'space-y-3 p-2' }, [
                  h(
                    'p',
                    { class: 'text-sm' },
                    'Tem certeza que deseja remover este usuário desta zona?',
                  ),
                  h(UButton, {
                    color: 'error',
                    variant: 'solid',
                    icon: 'i-lucide-check',
                    label: 'Confirmar',
                    onClick: () => {
                      deleteUser(row.original.zona, row.original.id)
                      close()
                    },
                  }),
                ]),
            }),
        ]),
    },
  ]
</script>

<template>
  <UModal
    v-model:open="modalUsers"
    title="Usuários"
    description="Gerencie os usuários associados a esta zona"
    :ui="{ footer: 'justify-end', content: 'max-w-4xl' }">
    <template #body>
      <div class="flex items-center justify-center space-x-3">
        <USelectMenu
          :disabled="isEditingUser"
          v-model="stateUser.email"
          label-key="nome"
          value-key="email"
          :items="members || []"
          icon="i-lucide-search"
          placeholder="Selecione o usuário..."
          class="mb-4">
          <template #item-label="{ item }">
            <div class="flex items-center gap-3">
              <img
                :src="`/server/api/file/${item.email}`"
                :alt="item.nome"
                class="size-6 cursor-pointer rounded-full" />
              <p>{{ item.nome }}</p>
            </div>
          </template>
        </USelectMenu>
        <USelect
          v-model="stateUser.permissao"
          :items="roleOptions"
          class="mb-4"
          placeholder="Selecione a permissão..."
          icon="i-lucide-shield-check" />
      </div>

      <UInput
        v-model="globalFilterUsers"
        class="mt-10 mb-4"
        placeholder="Buscar usuário..."
        icon="i-lucide-search" />

      <ClientOnly>
        <UTable
          ref="tableUsers"
          v-model:global-filter="globalFilterUsers"
          v-model:pagination="paginationUsers"
          class="mb-10"
          :pagination-options="{ getPaginationRowModel: getPaginationRowModel() }"
          :data="userData"
          :columns="columnsUsers" />
      </ClientOnly>

      <div
        v-if="userData && userData.length > paginationUsers.pageSize"
        class="flex justify-center border-t border-default pt-4">
        <UPagination
          active-
          active-variant="subtle"
          :default-page="(tableUsers?.tableApi?.getState().pagination.pageIndex || 0) + 1"
          :items-per-page="tableUsers?.tableApi?.getState().pagination.pageSize"
          :total="tableUsers?.tableApi?.getFilteredRowModel().rows.length"
          @update:page="(p) => tableUsers?.tableApi?.setPageIndex(p - 1)" />
      </div>
    </template>

    <template #footer>
      <UButton label="Fechar" :loading="isLoading" variant="outline" @click="modalUsers = false" />
      <UButton
        :label="stateUser.id ? 'Editar' : 'Adicionar'"
        icon="i-lucide-user-plus"
        :loading="isLoading"
        variant="outline"
        @click="addUser" />
    </template>
  </UModal>
</template>
