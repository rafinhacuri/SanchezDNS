<script setup lang="ts">
  import { safeParse } from 'valibot'

  import { NuxtTime } from '#components'

  useHead({ title: 'Solicitações' })

  const toast = useToast()

  const { isLoading, start, finish } = useLoadingIndicator()

  const page = ref(1)
  const itemsPerPage = ref(9)
  const filter = ref('')
  const filterDebounced = refDebounced(filter, 300)

  const { data, refresh } = await useApi<SolicitacaoResponse>('/solicitacoes', {
    method: 'GET',
    query: { page, limit: itemsPerPage, filter: filterDebounced },
    default: () => ({ solicitacoes: [], total: 0 }),
    transform: (response) => ({
      ...response,
      solicitacoes: response?.solicitacoes ?? [],
      total: response?.total ?? 0,
    }),
  })

  watch(filterDebounced, () => (page.value = 1))

  const totalSolicitacoes = computed(() => data.value.total)
  const hasSolicitacoes = computed(() => (data.value.solicitacoes ?? []).length > 0)

  const pendentSolicitacoes = computed(() =>
    (data.value.solicitacoes ?? []).filter((s) => s.status === 'pendente'),
  )

  const approvedSolicitacoes = computed(() =>
    (data.value.solicitacoes ?? []).filter((s) => s.status === 'aprovada'),
  )

  const rejectedSolicitacoes = computed(() =>
    (data.value.solicitacoes ?? []).filter((s) => s.status === 'rejeitada'),
  )

  type SolicitacaoAction = 'approve' | 'reject'

  const slideover = ref(false)
  const selectedSolicitacao = ref<Solicitacao | null>(null)
  const selectedAction = ref<SolicitacaoAction>('approve')

  const actionTitle = computed(() =>
    selectedAction.value === 'approve' ? 'Aprovar solicitação' : 'Rejeitar solicitação',
  )
  const actionDescription = computed(() =>
    selectedAction.value === 'approve'
      ? 'Confirme se deseja liberar o acesso deste usuário ao sistema.'
      : 'Confirme se deseja recusar esta solicitação de acesso.',
  )
  const actionButtonLabel = computed(() =>
    selectedAction.value === 'approve' ? 'Aprovar solicitação' : 'Rejeitar solicitação',
  )
  const actionIcon = computed(() =>
    selectedAction.value === 'approve' ? 'i-lucide-check-circle' : 'i-lucide-x-circle',
  )

  function openSolicitacaoSlideover(solicitacao: Solicitacao, action: SolicitacaoAction): void {
    selectedSolicitacao.value = solicitacao
    selectedAction.value = action
    slideover.value = true
  }

  function getStatusColor(st: string): 'neutral' | 'error' | 'warning' {
    switch (st) {
      case 'rejeitada': {
        return 'error'
      }
      case 'pendente': {
        return 'warning'
      }
      default: {
        return 'neutral'
      }
    }
  }

  async function status(): Promise<void> {
    start()

    const body = safeParse(StatusSchema, {
      id: selectedSolicitacao.value?.id ?? '',
      status: selectedAction.value === 'approve' ? 'aprovada' : 'rejeitada',
    })
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/ solicitacoes-status', {
      method: 'PUT',
      body: body.output,
    }).catch((error) => {
      toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
    })

    if (!res) return finish({ error: true })

    refresh()
    toast.add({ title: res.message, icon: 'i-lucide-shield-check', color: 'success' })
    slideover.value = false
    finish()
  }

  watch(slideover, (open) => {
    if (!open) selectedSolicitacao.value = null
  })
</script>

<template>
  <UContainer class="py-8">
    <div
      class="overflow-hidden rounded-3xl border border-default bg-linear-to-br from-primary/10 via-white to-white shadow-sm dark:via-gray-950 dark:to-gray-950">
      <div class="flex flex-col gap-6 p-6 sm:p-8 lg:flex-row lg:items-center lg:justify-between">
        <div class="max-w-2xl space-y-3">
          <UBadge color="primary" variant="soft" class="w-fit">Central de solicitações</UBadge>
          <div class="space-y-2">
            <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl dark:text-white">
              Solicitações
            </h1>
            <p class="text-base text-gray-600 dark:text-gray-400">
              Analise os pedidos pendentes, confira os dados do usuário e aprove ou rejeite
              solicitações rapidamente.
            </p>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3 sm:grid-cols-2 lg:min-w-80">
          <UCard variant="subtle" class="bg-white/70 dark:bg-gray-900/70">
            <div class="flex items-center gap-3">
              <div
                class="flex size-11 items-center justify-center rounded-2xl bg-primary/10 text-primary">
                <UIcon name="i-lucide-inbox" class="size-5" />
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Total</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">
                  {{ totalSolicitacoes }}
                </p>
              </div>
            </div>
          </UCard>
          <UCard variant="subtle" class="bg-white/70 dark:bg-gray-900/70">
            <div class="flex items-center gap-3">
              <div
                class="flex size-11 items-center justify-center rounded-2xl bg-success/10 text-success">
                <UIcon name="i-lucide-check-circle" class="size-5" />
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Aprovadas</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">
                  {{ approvedSolicitacoes.length }}
                </p>
              </div>
            </div>
          </UCard>

          <UCard variant="subtle" class="bg-white/70 dark:bg-gray-900/70">
            <div class="flex items-center gap-3">
              <div
                class="flex size-11 items-center justify-center rounded-2xl bg-error/10 text-error">
                <UIcon name="i-lucide-x-circle" class="size-5" />
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Rejeitadas</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">
                  {{ rejectedSolicitacoes.length }}
                </p>
              </div>
            </div>
          </UCard>

          <UCard variant="subtle" class="bg-white/70 dark:bg-gray-900/70">
            <div class="flex items-center gap-3">
              <div
                class="flex size-11 items-center justify-center rounded-2xl bg-warning/10 text-warning">
                <UIcon name="i-lucide-clock" class="size-5" />
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Pendentes</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">
                  {{ pendentSolicitacoes.length }}
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

    <div v-if="hasSolicitacoes" class="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <UCard
        v-for="solicitacao in data.solicitacoes ?? []"
        :key="solicitacao.email"
        class="group overflow-hidden transition hover:-translate-y-1 hover:shadow-lg">
        <div class="space-y-5">
          <div class="flex items-start justify-between gap-4">
            <div class="flex min-w-0 items-center gap-3">
              <img
                :src="`/server/api/file/${solicitacao.foto}`"
                :alt="solicitacao.email"
                class="size-12 rounded-2xl object-cover ring-2 ring-primary/10" />
              <div class="min-w-0">
                <h3 class="truncate text-base font-semibold text-gray-900 dark:text-white">
                  {{ solicitacao.nome }}
                </h3>
                <p class="truncate text-sm text-gray-500 dark:text-gray-400">
                  {{ solicitacao.email }}
                </p>
              </div>
            </div>

            <UBadge :color="getStatusColor(solicitacao.status)" variant="soft">
              {{ solicitacao.status }}
            </UBadge>
          </div>

          <div class="rounded-2xl bg-gray-50 p-4 dark:bg-gray-900/60">
            <div class="flex items-center gap-3 text-sm text-gray-600 dark:text-gray-400">
              <UIcon name="i-lucide-calendar-clock" class="size-5 text-primary" />
              <span>Solicitado em</span>
            </div>
            <NuxtTime
              :datetime="solicitacao.createdAt"
              day="2-digit"
              month="2-digit"
              year="2-digit"
              hour="2-digit"
              minute="2-digit"
              locale="pt-BR"
              class="mt-1 block text-lg font-semibold text-gray-900 dark:text-white" />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <UButton
              :loading="isLoading"
              v-if="solicitacao.status === 'pendente'"
              label="Rejeitar"
              color="error"
              variant="soft"
              icon="i-lucide-x"
              block
              @click="openSolicitacaoSlideover(solicitacao, 'reject')" />
            <UButton
              :loading="isLoading"
              v-if="solicitacao.status === 'pendente' || solicitacao.status === 'rejeitada'"
              label="Aprovar"
              color="success"
              icon="i-lucide-check"
              block
              @click="openSolicitacaoSlideover(solicitacao, 'approve')" />
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
    <USlideover v-model:open="slideover" :title="actionTitle" :description="actionDescription">
      <template #body>
        <div v-if="selectedSolicitacao" class="space-y-6">
          <div class="rounded-3xl border border-default bg-gray-50 p-5 dark:bg-gray-900/60">
            <div class="flex items-center gap-4">
              <img
                :src="`/server/api/file/${selectedSolicitacao.foto}`"
                :alt="selectedSolicitacao.email"
                class="size-16 rounded-2xl object-cover ring-2 ring-primary/10" />
              <div class="min-w-0">
                <h3 class="truncate text-lg font-semibold text-gray-900 dark:text-white">
                  {{ selectedSolicitacao.nome }}
                </h3>
                <p class="truncate text-sm text-gray-500 dark:text-gray-400">
                  {{ selectedSolicitacao.email }}
                </p>
              </div>
            </div>
          </div>

          <div class="space-y-3 rounded-3xl border border-default p-5">
            <div class="flex items-center gap-3">
              <div
                class="flex size-11 items-center justify-center rounded-2xl"
                :class="
                  selectedAction === 'approve'
                    ? 'bg-success/10 text-success'
                    : 'bg-error/10 text-error'
                ">
                <UIcon :name="actionIcon" class="size-5" />
              </div>
              <div>
                <p class="font-semibold text-gray-900 dark:text-white">{{ actionTitle }}</p>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  Esta ação precisa ser confirmada antes de continuar.
                </p>
              </div>
            </div>

            <UAlert
              v-if="selectedAction === 'approve'"
              color="info"
              variant="soft"
              icon="i-lucide-info"
              title="O que acontece após aprovar?"
              description="Ao aprovar, o cadastro da pessoa será criado com level member e ela poderá fazer login. Para tornar a pessoa administradora, entre na página Cadastros e edite o level do usuário para admin, liberando acesso total a todas as zonas para adicionar, editar ou remover registros. Para liberar acesso apenas a zonas específicas, um administrador deve entrar na zona desejada e adicionar o usuário com escrita ou leitura naquela zona." />

            <UAlert
              v-else
              color="error"
              variant="soft"
              icon="i-lucide-triangle-alert"
              title="Confirmação de recusa"
              description="Ao recusar, esta solicitação não criará cadastro para o usuário e ele não poderá acessar o sistema com esses dados." />

            <div class="rounded-2xl bg-gray-50 p-4 dark:bg-gray-900/60">
              <div class="flex items-center gap-3 text-sm text-gray-600 dark:text-gray-400">
                <UIcon name="i-lucide-calendar-clock" class="size-5 text-primary" />
                <span>Solicitado em</span>
              </div>
              <NuxtTime
                :datetime="selectedSolicitacao.createdAt"
                day="2-digit"
                month="2-digit"
                year="2-digit"
                hour="2-digit"
                minute="2-digit"
                locale="pt-BR"
                class="mt-1 block text-base font-semibold text-gray-900 dark:text-white" />
            </div>
          </div>
        </div>
      </template>

      <template #footer>
        <div class="grid w-full grid-cols-2 gap-3">
          <UButton
            :loading="isLoading"
            label="Cancelar"
            color="neutral"
            variant="soft"
            block
            @click="
              () => {
                slideover = false
              }
            " />
          <UButton
            @click="status"
            :loading="isLoading"
            :label="actionButtonLabel"
            :color="selectedAction === 'approve' ? 'success' : 'error'"
            :icon="selectedAction === 'approve' ? 'i-lucide-check' : 'i-lucide-x'"
            block />
        </div>
      </template>
    </USlideover>
  </UContainer>
</template>
