<script setup lang="ts">
  import { safeParse } from 'valibot'

  import { NuxtTime } from '#components'

  useHead({ title: 'Cadastros' })

  const toast = useToast()

  const { isLoading, start, finish } = useLoadingIndicator()

  const page = ref(1)
  const itemsPerPage = ref(9)
  const filter = ref('')
  const filterDebounced = refDebounced(filter, 300)

  const { data } = await useFetch<CadastroResponse>('/server/api/cadastros', {
    method: 'GET',
    query: { page, limit: itemsPerPage, filter: filterDebounced },
    default: () => ({ cadastros: [], total: 0 }),
    transform: (response) => ({
      ...response,
      cadastros: response?.cadastros ?? [],
      total: response?.total ?? 0,
    }),
  })

  watch(filterDebounced, () => (page.value = 1))

  const totalCadastros = computed(() => data.value.total)
  const hasCadastros = computed(() => (data.value.cadastros ?? []).length > 0)

  function getStatusColor(st: string): 'neutral' | 'error' | 'warning' {
    switch (st) {
      case 'admin': {
        return 'warning'
      }
      default: {
        return 'neutral'
      }
    }
  }
</script>

<template>
  <UContainer class="py-8">
    <div
      class="overflow-hidden rounded-3xl border border-default bg-linear-to-br from-primary/10 via-white to-white shadow-sm dark:via-gray-950 dark:to-gray-950">
      <div class="flex flex-col gap-6 p-6 sm:p-8 lg:flex-row lg:items-center lg:justify-between">
        <div class="max-w-2xl space-y-3">
          <UBadge color="primary" variant="soft" class="w-fit">Central de Cadastros</UBadge>
          <div class="space-y-2">
            <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl dark:text-white">
              Cadastros
            </h1>
            <p class="text-base text-gray-600 dark:text-gray-400">
              Gerencie os cadastros de usuários que possuem acesso ao sistema. Aqui você pode
              visualizar, editar ou remover cadastros existentes.
            </p>
          </div>
        </div>

        <div class="lg:min-w-80">
          <UCard variant="subtle" class="bg-white/70 dark:bg-gray-900/70">
            <div class="flex items-center gap-3">
              <div
                class="flex size-11 items-center justify-center rounded-2xl bg-primary/10 text-primary">
                <UIcon name="i-lucide-inbox" class="size-5" />
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Total</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">
                  {{ totalCadastros }}
                </p>
              </div>
            </div>
          </UCard>
        </div>
      </div>
    </div>

    <div class="mt-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Pedidos recentes</h2>
        <p class="text-sm text-gray-500 dark:text-gray-400">
          Gerencie as solicitações recebidas pelo sistema.
        </p>
      </div>

      <UInput
        v-model="filter"
        icon="i-lucide-search"
        placeholder="Pesquisar por nome ou email..."
        class="sm:w-80" />
    </div>

    <div v-if="hasCadastros" class="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <UCard
        v-for="cadastro in data.cadastros ?? []"
        :key="cadastro.email"
        class="group overflow-hidden transition hover:-translate-y-1 hover:shadow-lg">
        <div class="space-y-5">
          <div class="flex items-start justify-between gap-4">
            <div class="flex min-w-0 items-center gap-3">
              <img
                :src="`/server/api/file/${cadastro.foto}`"
                :alt="cadastro.email"
                class="size-12 rounded-2xl object-cover ring-2 ring-primary/10" />
              <div class="min-w-0">
                <h3 class="truncate text-base font-semibold text-gray-900 dark:text-white">
                  {{ cadastro.nome }}
                </h3>
                <p class="truncate text-sm text-gray-500 dark:text-gray-400">
                  {{ cadastro.email }}
                </p>
              </div>
            </div>

            <UBadge :color="getStatusColor(cadastro.level)" variant="soft">
              {{ cadastro.level === 'admin' ? 'Administrador' : 'Usuário' }}
            </UBadge>
          </div>

          <div class="rounded-2xl bg-gray-50 p-4 dark:bg-gray-900/60">
            <div class="flex items-center gap-3 text-sm text-gray-600 dark:text-gray-400">
              <UIcon name="i-lucide-calendar-clock" class="size-5 text-primary" />
              <span>Criado em</span>
            </div>
            <NuxtTime
              :datetime="cadastro.createdAt"
              day="2-digit"
              month="2-digit"
              year="2-digit"
              hour="2-digit"
              minute="2-digit"
              locale="pt-BR"
              class="mt-1 block text-lg font-semibold text-gray-900 dark:text-white" />
          </div>

          <div class="space-y-3">
            <UButton
              :loading="isLoading"
              label="Editar"
              color="info"
              variant="soft"
              icon="i-lucide-pen"
              block />
            <UButton
              :loading="isLoading"
              label="Alterar Level"
              color="warning"
              icon="i-lucide-user-pen"
              block
              variant="soft" />
            <UButton
              :loading="isLoading"
              label="Excluir"
              color="error"
              icon="i-lucide-trash-2"
              variant="soft"
              block />
          </div>
        </div>
      </UCard>
    </div>

    <UCard v-else class="mt-6 overflow-hidden border-dashed">
      <div class="flex flex-col items-center justify-center py-12 text-center">
        <div
          class="mb-4 flex size-20 items-center justify-center rounded-3xl bg-primary/10 text-primary">
          <UIcon name="i-lucide-inbox" class="size-10" />
        </div>
        <h3 class="text-xl font-semibold text-gray-900 dark:text-white">
          Nenhuma solicitação pendente
        </h3>
        <p class="mt-2 max-w-md text-sm text-gray-500 dark:text-gray-400">
          Não existem solicitações aguardando aprovação no momento. Quando novos usuários
          solicitarem acesso ao sistema, eles aparecerão aqui.
        </p>
      </div>
    </UCard>

    <div
      v-if="data.total > itemsPerPage"
      class="mt-8 flex justify-center border-t border-default pt-6">
      <UPagination
        v-model:page="page"
        active-variant="subtle"
        :total="data.total"
        :items-per-page="itemsPerPage" />
    </div>
  </UContainer>
</template>
