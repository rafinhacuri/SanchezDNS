<script setup lang="ts">
  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()
  const { copy } = useClipboard()

  const props = defineProps({
    zone: { type: String, required: true },
    isAdmin: { type: Boolean, default: false },
  })

  const zoneId = computed(() => props.zone)

  const modal = defineModel<boolean>('open', { default: false })

  const status = ref<DnssecStatus>()
  const carregando = ref(false)

  async function carregar(): Promise<void> {
    carregando.value = true

    const res = await $api<DnssecStatus>('/zone/dnssec', {
      method: 'GET',
      query: { zone: zoneId.value },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao buscar o DNSSEC da zona',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    carregando.value = false

    if (res) status.value = res
  }

  watch(modal, (nv) => {
    if (nv) void carregar()
    else status.value = undefined
  })

  const dsMeta = computed(() => {
    const estado = status.value?.dsStatus

    if (estado === 'ok') {
      return {
        label: 'DS publicado no pai',
        detail:
          'A cadeia de confiança está fechada: resolvers validantes rejeitam respostas forjadas para esta zona.',
        icon: 'i-lucide-shield-check',
        tone: 'text-emerald-600 dark:text-emerald-400',
        box: 'border-emerald-500/20 bg-emerald-500/5',
      }
    }

    if (estado === 'divergente') {
      return {
        label: 'DS do pai não confere',
        detail:
          'O pai publica um DS que não corresponde a nenhuma chave atual da zona. Resolvers validantes deixam de resolver o domínio até o DS ser corrigido no registrador.',
        icon: 'i-lucide-shield-x',
        tone: 'text-rose-600 dark:text-rose-400',
        box: 'border-rose-500/20 bg-rose-500/5',
      }
    }

    if (estado === 'ausente') {
      return {
        label: 'Sem DS no pai',
        detail:
          'A zona está assinada, mas o registrador não publicou o DS. Sem ele nenhum resolver valida as assinaturas, e a zona continua exposta a envenenamento de cache. Publique o DS abaixo no registrador.',
        icon: 'i-lucide-shield-alert',
        tone: 'text-amber-600 dark:text-amber-400',
        box: 'border-amber-500/20 bg-amber-500/5',
      }
    }

    return {
      label: 'Pai indisponível',
      detail:
        'Não foi possível consultar o DS no pai através dos resolvers públicos. Tente novamente em instantes.',
      icon: 'i-lucide-shield-question',
      tone: 'text-dimmed',
      box: 'border-default bg-muted/40',
    }
  })

  const dsRecords = computed(() =>
    (status.value?.keys ?? []).filter((key) => key.ds?.length).flatMap((key) => key.ds),
  )

  async function aplicarNsec3(): Promise<void> {
    start()

    const res = await $api('/zone/nsec3', {
      method: 'PATCH',
      query: { zone: zoneId.value },
    }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao aplicar o NSEC3',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    finish({ force: true })
    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    await carregar()
  }

  function copiar(valor: string): void {
    copy(valor)
    toast.add({
      title: 'DS copiado para a área de transferência',
      color: 'success',
      icon: 'i-lucide-circle-check',
    })
  }
</script>

<template>
  <UModal
    v-model:open="modal"
    title="DNSSEC"
    :description="`Estado da assinatura e da cadeia de confiança de ${zoneId}`"
    :ui="{ content: 'max-w-2xl', footer: 'justify-end' }">
    <template #body>
      <div v-if="carregando && !status" class="flex items-center justify-center gap-2.5 py-12">
        <UIcon name="i-lucide-loader-circle" class="size-4 animate-spin text-dimmed" />
        <p class="text-sm text-muted">Consultando o servidor e o pai da zona...</p>
      </div>

      <div v-else-if="status" class="space-y-4">
        <div class="grid gap-3 sm:grid-cols-2">
          <div class="rounded-lg border border-default bg-muted/40 p-3">
            <p class="text-[11px] font-semibold tracking-wider text-dimmed uppercase">Assinatura</p>
            <div class="mt-1.5 flex items-center gap-1.5">
              <span
                class="size-1.5 rounded-full"
                :class="status.dnssec ? 'bg-emerald-500' : 'bg-slate-400 dark:bg-slate-600'" />
              <span class="text-sm" :class="status.dnssec ? 'text-toned' : 'text-dimmed'">
                {{ status.dnssec ? 'Zona assinada' : 'Zona sem DNSSEC' }}
              </span>
            </div>
          </div>

          <div class="rounded-lg border border-default bg-muted/40 p-3">
            <p class="text-[11px] font-semibold tracking-wider text-dimmed uppercase">
              Prova de negação
            </p>
            <div class="mt-1.5 flex items-center gap-1.5">
              <span
                class="size-1.5 rounded-full"
                :class="status.nsec3 ? 'bg-emerald-500' : 'bg-amber-500'" />
              <span
                class="data text-sm"
                :class="status.nsec3 ? 'text-toned' : 'text-amber-600 dark:text-amber-400'">
                {{ status.nsec3 ? 'NSEC3' : 'NSEC' }}
              </span>
            </div>
          </div>
        </div>

        <div
          v-if="status.dnssec && !status.nsec3"
          class="flex gap-3 rounded-lg border border-amber-500/20 bg-amber-500/5 p-3">
          <UIcon name="i-lucide-triangle-alert" class="mt-0.5 size-4 shrink-0 text-amber-500" />
          <div class="min-w-0 space-y-2">
            <p class="text-sm text-toned">
              Esta zona usa NSEC, que permite enumerar todos os registros com ferramentas como
              <span class="data">ldns-walk</span>. Migrar para NSEC3 fecha esse reconhecimento.
            </p>
            <UButton
              v-if="isAdmin"
              size="xs"
              color="warning"
              icon="i-lucide-shield-plus"
              label="Aplicar NSEC3"
              :loading="isLoading"
              @click="aplicarNsec3" />
          </div>
        </div>

        <div class="rounded-lg border p-3" :class="dsMeta.box">
          <div class="flex gap-3">
            <UIcon :name="dsMeta.icon" class="mt-0.5 size-4 shrink-0" :class="dsMeta.tone" />
            <div class="min-w-0 space-y-1">
              <p class="text-sm font-semibold" :class="dsMeta.tone">{{ dsMeta.label }}</p>
              <p class="text-sm text-muted">{{ dsMeta.detail }}</p>
              <p v-if="status.dsStatus === 'ok' && status.validado" class="text-xs text-dimmed">
                Resposta validada pelo resolver (flag AD).
              </p>
            </div>
          </div>
        </div>

        <div v-if="dsRecords.length" class="space-y-2">
          <div class="flex items-center gap-2">
            <UIcon name="i-lucide-key-round" class="size-4 text-dimmed" />
            <h3 class="text-sm font-semibold text-highlighted">DS para publicar no registrador</h3>
          </div>

          <div
            v-for="ds in dsRecords"
            :key="ds"
            class="flex items-start gap-2 rounded-lg border border-default bg-muted/40 p-2.5">
            <p class="min-w-0 flex-1 data text-xs break-all text-toned">{{ ds }}</p>
            <UButton
              size="xs"
              color="neutral"
              variant="ghost"
              icon="i-lucide-copy"
              aria-label="Copiar DS"
              @click="copiar(ds)" />
          </div>
        </div>

        <div v-if="status.dsPai.length" class="space-y-2">
          <div class="flex items-center gap-2">
            <UIcon name="i-lucide-globe" class="size-4 text-dimmed" />
            <h3 class="text-sm font-semibold text-highlighted">DS encontrado no pai</h3>
          </div>

          <div
            v-for="ds in status.dsPai"
            :key="ds"
            class="rounded-lg border border-default bg-muted/40 p-2.5">
            <p class="data text-xs break-all text-muted">{{ ds }}</p>
          </div>
        </div>
      </div>
    </template>

    <template #footer>
      <UButton
        label="Fechar"
        color="neutral"
        variant="ghost"
        @click="
          () => {
            modal = false
          }
        " />
      <UButton
        label="Reconsultar"
        icon="i-lucide-refresh-cw"
        variant="outline"
        color="neutral"
        :loading="carregando"
        @click="carregar" />
    </template>
  </UModal>
</template>
