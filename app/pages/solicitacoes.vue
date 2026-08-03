<script setup lang="ts">
  import { safeParse } from 'valibot'

  import { NuxtTime } from '#components'

  useHead({ title: 'Solicitações' })

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()
  const { statusSchema } = useCadastroSchema()

  const baseUrl = useApiUrl()

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

  const solicitacoes = computed(() => data.value.solicitacoes ?? [])
  const hasSolicitacoes = computed(() => solicitacoes.value.length > 0)

  const kpis = computed(() => [
    {
      label: 'Total',
      value: data.value.total,
      icon: 'i-lucide-inbox',
      tint: 'text-blue-500',
    },
    {
      label: 'Pendentes',
      value: solicitacoes.value.filter((s) => s.status === 'pendente').length,
      icon: 'i-lucide-clock',
      tint: 'text-amber-500',
    },
    {
      label: 'Aprovadas',
      value: solicitacoes.value.filter((s) => s.status === 'aprovada').length,
      icon: 'i-lucide-circle-check',
      tint: 'text-emerald-500',
    },
    {
      label: 'Rejeitadas',
      value: solicitacoes.value.filter((s) => s.status === 'rejeitada').length,
      icon: 'i-lucide-circle-x',
      tint: 'text-rose-500',
    },
  ])

  function statusChip(st: string): string {
    if (st === 'aprovada') {
      return 'bg-emerald-500/10 text-emerald-600 ring-emerald-500/25 dark:text-emerald-400'
    }
    if (st === 'rejeitada') {
      return 'bg-rose-500/10 text-rose-600 ring-rose-500/25 dark:text-rose-400'
    }
    return 'bg-amber-500/10 text-amber-600 ring-amber-500/25 dark:text-amber-400'
  }

  type SolicitacaoAction = 'approve' | 'reject'

  const slideover = ref(false)
  const selectedSolicitacao = ref<Solicitacao | null>(null)
  const selectedAction = ref<SolicitacaoAction>('approve')

  const isApprove = computed(() => selectedAction.value === 'approve')

  const actionTitle = computed(() =>
    isApprove.value ? 'Aprovar solicitação' : 'Rejeitar solicitação',
  )

  const actionDescription = computed(() =>
    isApprove.value
      ? 'Confirme se deseja liberar o acesso deste usuário ao sistema.'
      : 'Confirme se deseja recusar esta solicitação de acesso.',
  )

  function openSolicitacaoSlideover(solicitacao: Solicitacao, action: SolicitacaoAction): void {
    selectedSolicitacao.value = solicitacao
    selectedAction.value = action
    slideover.value = true
  }

  async function status(): Promise<void> {
    start()

    const body = safeParse(statusSchema, {
      id: selectedSolicitacao.value?.id ?? '',
      status: isApprove.value ? 'aprovada' : 'rejeitada',
    })

    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/solicitacoes/status', {
      method: 'PUT',
      body: body.output,
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao atualizar a solicitação',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    await refresh()
    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    slideover.value = false
    finish()
  }

  watch(slideover, (open) => {
    if (!open) selectedSolicitacao.value = null
  })
</script>

<template>
  <UContainer class="py-8 sm:py-10">
    <PageHeader
      eyebrow="Controle de acesso"
      title="Solicitações"
      description="Analise os pedidos de acesso, confira os dados do usuário e aprove ou rejeite.">
      <template #actions>
        <UInput
          v-model="filter"
          icon="i-lucide-search"
          placeholder="Pesquisar por nome ou email..."
          class="w-full sm:w-80"
          :ui="{ base: 'data' }" />
      </template>
    </PageHeader>

    <div class="stagger mb-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="kpi of kpis"
        :key="kpi.label"
        class="console-rail overflow-hidden rounded-xl border border-default bg-default p-5">
        <div class="flex items-start justify-between gap-4">
          <div class="min-w-0 space-y-1">
            <p class="text-xs font-semibold tracking-wider text-dimmed uppercase">
              {{ kpi.label }}
            </p>
            <p class="text-3xl font-bold text-highlighted tnum">{{ kpi.value }}</p>
          </div>
          <UIcon :name="kpi.icon" class="size-5 shrink-0" :class="kpi.tint" />
        </div>
      </div>
    </div>

    <div v-if="hasSolicitacoes" class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <div
        v-for="solicitacao of solicitacoes"
        :key="solicitacao.email"
        class="flex flex-col gap-4 rounded-xl border border-default bg-default p-5 transition-colors duration-150 hover:border-accented">
        <div class="flex items-start justify-between gap-3">
          <div class="flex min-w-0 items-center gap-3">
            <img
              :src="`${baseUrl}/file/${solicitacao.foto}`"
              :alt="solicitacao.email"
              class="size-11 shrink-0 rounded-full object-cover ring-1 ring-default" />
            <div class="min-w-0">
              <p class="truncate text-sm font-semibold text-highlighted">
                {{ solicitacao.nome }}
              </p>
              <p class="truncate data text-xs text-dimmed">{{ solicitacao.email }}</p>
            </div>
          </div>

          <span
            class="shrink-0 rounded-md px-2 py-0.5 text-[11px] font-semibold tracking-wide uppercase ring-1 ring-inset"
            :class="statusChip(solicitacao.status)">
            {{ solicitacao.status }}
          </span>
        </div>

        <div class="flex items-center gap-2 border-t border-default pt-3 text-xs text-muted">
          <UIcon name="i-lucide-calendar-clock" class="size-4 shrink-0 text-dimmed" />
          <span>Solicitado em</span>
          <NuxtTime
            :datetime="solicitacao.createdAt"
            day="2-digit"
            month="2-digit"
            year="2-digit"
            hour="2-digit"
            minute="2-digit"
            locale="pt-BR"
            class="data font-medium text-toned tnum" />
        </div>

        <div
          v-if="solicitacao.status !== 'aprovada'"
          class="mt-auto grid gap-2"
          :class="solicitacao.status === 'pendente' ? 'grid-cols-2' : 'grid-cols-1'">
          <UButton
            v-if="solicitacao.status === 'pendente'"
            label="Rejeitar"
            color="error"
            variant="soft"
            icon="i-lucide-x"
            block
            :loading="isLoading"
            @click="openSolicitacaoSlideover(solicitacao, 'reject')" />
          <UButton
            label="Aprovar"
            color="success"
            icon="i-lucide-check"
            block
            :loading="isLoading"
            @click="openSolicitacaoSlideover(solicitacao, 'approve')" />
        </div>
      </div>
    </div>

    <div
      v-else
      class="flex flex-col items-center gap-2 rounded-xl border border-default bg-default py-16 text-center">
      <UIcon name="i-lucide-inbox" class="size-7 text-dimmed" />
      <p class="text-sm font-medium text-toned">Nenhuma solicitação encontrada</p>
      <p class="max-w-sm text-xs text-dimmed">
        {{
          filter
            ? 'Nenhum resultado para esse filtro. Tente outro termo.'
            : 'Quando alguém solicitar acesso ao sistema, o pedido aparece aqui.'
        }}
      </p>
    </div>

    <div v-if="data.total > itemsPerPage" class="mt-6 flex justify-center">
      <UPagination
        v-model:page="page"
        active-variant="subtle"
        :total="data.total"
        :items-per-page="itemsPerPage" />
    </div>

    <USlideover
      v-model:open="slideover"
      :title="actionTitle"
      :description="actionDescription"
      :ui="{ footer: 'justify-end' }">
      <template #body>
        <div v-if="selectedSolicitacao" class="space-y-5">
          <div class="flex items-center gap-4 rounded-lg border border-default bg-muted/40 p-4">
            <img
              :src="`${baseUrl}/file/${selectedSolicitacao.foto}`"
              :alt="selectedSolicitacao.email"
              class="size-14 shrink-0 rounded-full object-cover ring-1 ring-default" />
            <div class="min-w-0">
              <p class="truncate text-base font-semibold text-highlighted">
                {{ selectedSolicitacao.nome }}
              </p>
              <p class="truncate data text-sm text-dimmed">{{ selectedSolicitacao.email }}</p>
            </div>
          </div>

          <UAlert
            v-if="isApprove"
            color="info"
            variant="soft"
            icon="i-lucide-info"
            title="O que acontece após aprovar?"
            description="O cadastro é criado com level member e a pessoa já consegue fazer login. Para torná-la administradora, mude o level em Cadastros. Para liberar zonas específicas, adicione o usuário com leitura ou escrita dentro da zona." />

          <UAlert
            v-else
            color="error"
            variant="soft"
            icon="i-lucide-triangle-alert"
            title="Confirmação de recusa"
            description="A solicitação não vai gerar cadastro e a pessoa não conseguirá acessar o sistema com esses dados." />

          <div class="flex items-center gap-2 text-xs text-muted">
            <UIcon name="i-lucide-calendar-clock" class="size-4 shrink-0 text-dimmed" />
            <span>Solicitado em</span>
            <NuxtTime
              :datetime="selectedSolicitacao.createdAt"
              day="2-digit"
              month="2-digit"
              year="2-digit"
              hour="2-digit"
              minute="2-digit"
              locale="pt-BR"
              class="data font-medium text-toned tnum" />
          </div>
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
              slideover = false
            }
          " />
        <UButton
          :label="actionTitle"
          :color="isApprove ? 'success' : 'error'"
          :icon="isApprove ? 'i-lucide-check' : 'i-lucide-x'"
          :loading="isLoading"
          @click="status" />
      </template>
    </USlideover>
  </UContainer>
</template>
