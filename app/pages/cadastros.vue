<script setup lang="ts">
  import { safeParse } from 'valibot'

  import { NuxtTime } from '#components'

  useHead({ title: 'Cadastros' })

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()
  const { idSchema, cadastroEditSchema } = useCadastroSchema()

  const baseUrl = useApiUrl()

  const page = ref(1)
  const itemsPerPage = ref(9)
  const filter = ref('')
  const filterDebounced = refDebounced(filter, 300)

  const { data, refresh } = await useApi<CadastroResponse>('/cadastros', {
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

  const cadastros = computed(() => data.value.cadastros ?? [])
  const hasCadastros = computed(() => cadastros.value.length > 0)

  const kpis = computed(() => [
    {
      label: 'Total',
      value: data.value.total,
      icon: 'i-lucide-users',
      tint: 'text-blue-500',
    },
    {
      label: 'Administradores',
      value: cadastros.value.filter((c) => c.level === 'admin').length,
      icon: 'i-lucide-shield',
      tint: 'text-amber-500',
    },
    {
      label: 'Membros',
      value: cadastros.value.filter((c) => c.level !== 'admin').length,
      icon: 'i-lucide-user',
      tint: 'text-emerald-500',
    },
  ])

  const idDelete = ref('')
  const nomeDelete = ref('')
  const modalDelete = ref(false)

  function openDelete(id: string, nome: string): void {
    idDelete.value = id
    nomeDelete.value = nome
    modalDelete.value = true
  }

  async function deleteCadastro(): Promise<void> {
    start()

    const body = safeParse(idSchema, { id: idDelete.value })
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/cadastro', { method: 'DELETE', body: body.output }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao deletar o cadastro',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    await refresh()
    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    modalDelete.value = false
    finish()
  }

  watch(modalDelete, (open) => {
    if (!open) {
      idDelete.value = ''
      nomeDelete.value = ''
    }
  })

  const idLevel = ref('')
  const nomeLevel = ref('')
  const level = ref('')
  const modalLevel = ref(false)

  function openLevel(id: string, nome: string, lvl: string): void {
    idLevel.value = id
    nomeLevel.value = nome
    level.value = lvl
    modalLevel.value = true
  }

  const novoLevel = computed(() => (level.value === 'admin' ? 'member' : 'admin'))

  async function updateLevel(): Promise<void> {
    start()

    const body = safeParse(idSchema, { id: idLevel.value })
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/level', { method: 'PATCH', body: body.output }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao alterar o nível',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    await refresh()
    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    modalLevel.value = false
    finish()
  }

  watch(modalLevel, (open) => {
    if (!open) {
      idLevel.value = ''
      nomeLevel.value = ''
      level.value = ''
    }
  })

  const modalEdit = ref(false)
  const displayName = ref('')
  const editCadastro = ref<CadastroEditSchema>({ id: '', senha: '', foto: '', nome: '' })
  const confirmSenha = ref('')
  const foto = ref<File | null>(null)
  const show = ref(false)
  const showConfirm = ref(false)

  function openEdit(cadastro: Cadastro): void {
    editCadastro.value = { id: cadastro.id, senha: '', foto: cadastro.foto, nome: cadastro.nome }
    displayName.value = cadastro.nome
    modalEdit.value = true
  }

  function createObjectUrl(file: File): string {
    return URL.createObjectURL(file)
  }

  const requisitos = [
    { regex: /.{8,}/u, text: 'Pelo menos 8 caracteres' },
    { regex: /\d/u, text: 'Pelo menos 1 número' },
    { regex: /[a-z]/u, text: 'Pelo menos 1 letra minúscula' },
    { regex: /[A-Z]/u, text: 'Pelo menos 1 letra maiúscula' },
  ]

  const strength = computed(() =>
    requisitos.map((req) => ({ met: req.regex.test(editCadastro.value.senha), text: req.text })),
  )

  const score = computed(() => strength.value.filter((req) => req.met).length)

  const color = computed(() => {
    if (score.value === 0) return 'neutral'
    if (score.value <= 1) return 'error'
    if (score.value <= 3) return 'warning'
    return 'success'
  })

  const text = computed(() => {
    if (score.value === 0) return 'Deixe em branco para manter'
    if (score.value <= 2) return 'Senha fraca'
    if (score.value === 3) return 'Senha média'
    return 'Senha forte'
  })

  watch(foto, async (file) => {
    if (!file) return

    const formData = new FormData()
    formData.append('file', file)

    const res = await $api('/file', { method: 'PUT', body: formData }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao enviar a foto',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (res) editCadastro.value.foto = res.message
  })

  async function editar(): Promise<void> {
    start()

    if (editCadastro.value.senha !== confirmSenha.value) {
      toast.add({ title: 'As senhas não coincidem', icon: 'i-lucide-shield-alert', color: 'error' })
      return finish({ error: true })
    }

    const body = safeParse(cadastroEditSchema, editCadastro.value)
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/cadastro', { method: 'PUT', body: body.output }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao atualizar o cadastro',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
    await refresh()
    modalEdit.value = false
    finish()
  }

  watch(modalEdit, (open) => {
    if (!open) {
      editCadastro.value = { id: '', senha: '', foto: '', nome: '' }
      displayName.value = ''
      foto.value = null
      confirmSenha.value = ''
    }
  })
</script>

<template>
  <UContainer class="py-8 sm:py-10">
    <PageHeader
      eyebrow="Controle de acesso"
      title="Cadastros"
      description="Gerencie quem tem acesso ao sistema, o nível de cada um e os dados da conta.">
      <template #actions>
        <UInput
          v-model="filter"
          icon="i-lucide-search"
          placeholder="Pesquisar por nome ou email..."
          class="w-full sm:w-80"
          :ui="{ base: 'data' }" />
      </template>
    </PageHeader>

    <div class="stagger mb-6 grid gap-4 sm:grid-cols-3">
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

    <div v-if="hasCadastros" class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <div
        v-for="cadastro of cadastros"
        :key="cadastro.email"
        class="flex flex-col gap-4 rounded-xl border border-default bg-default p-5 transition-colors duration-150 hover:border-accented">
        <div class="flex items-start justify-between gap-3">
          <div class="flex min-w-0 items-center gap-3">
            <img
              :src="`${baseUrl}/file/${cadastro.foto}`"
              :alt="cadastro.email"
              class="size-11 shrink-0 rounded-full object-cover ring-1 ring-default" />
            <div class="min-w-0">
              <p class="truncate text-sm font-semibold text-highlighted">{{ cadastro.nome }}</p>
              <p class="truncate data text-xs text-dimmed">{{ cadastro.email }}</p>
            </div>
          </div>

          <span
            class="shrink-0 rounded-md px-2 py-0.5 text-[11px] font-semibold tracking-wide uppercase ring-1 ring-inset"
            :class="
              cadastro.level === 'admin'
                ? 'bg-amber-500/10 text-amber-600 ring-amber-500/25 dark:text-amber-400'
                : 'bg-slate-500/10 text-slate-600 ring-slate-500/25 dark:text-slate-300'
            ">
            {{ cadastro.level === 'admin' ? 'Admin' : 'Membro' }}
          </span>
        </div>

        <div class="flex items-center gap-2 border-t border-default pt-3 text-xs text-muted">
          <UIcon name="i-lucide-calendar-clock" class="size-4 shrink-0 text-dimmed" />
          <span>Criado em</span>
          <NuxtTime
            :datetime="cadastro.createdAt"
            day="2-digit"
            month="2-digit"
            year="2-digit"
            hour="2-digit"
            minute="2-digit"
            locale="pt-BR"
            class="data font-medium text-toned tnum" />
        </div>

        <div class="mt-auto grid grid-cols-3 gap-2">
          <UButton
            label="Editar"
            color="neutral"
            variant="outline"
            icon="i-lucide-pencil"
            size="sm"
            block
            :loading="isLoading"
            @click="openEdit(cadastro)" />
          <UButton
            :label="cadastro.level === 'admin' ? 'Membro' : 'Admin'"
            color="neutral"
            variant="outline"
            icon="i-lucide-shield"
            size="sm"
            block
            :loading="isLoading"
            @click="openLevel(cadastro.id, cadastro.nome, cadastro.level)" />
          <UButton
            label="Excluir"
            color="error"
            variant="soft"
            icon="i-lucide-trash-2"
            size="sm"
            block
            :loading="isLoading"
            @click="openDelete(cadastro.id, cadastro.nome)" />
        </div>
      </div>
    </div>

    <div
      v-else
      class="flex flex-col items-center gap-2 rounded-xl border border-default bg-default py-16 text-center">
      <UIcon name="i-lucide-users" class="size-7 text-dimmed" />
      <p class="text-sm font-medium text-toned">Nenhum cadastro encontrado</p>
      <p class="max-w-sm text-xs text-dimmed">
        {{
          filter
            ? 'Nenhum resultado para esse filtro. Tente outro termo.'
            : 'Assim que uma solicitação for aprovada, o cadastro aparece aqui.'
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

    <UModal
      v-model:open="modalDelete"
      title="Excluir cadastro"
      description="Esta ação é permanente e não pode ser desfeita."
      :ui="{ footer: 'justify-end' }">
      <template #body>
        <div
          class="flex gap-3 rounded-lg border border-rose-500/20 bg-rose-500/5 p-3 text-sm text-toned">
          <UIcon name="i-lucide-triangle-alert" class="mt-0.5 size-4 shrink-0 text-rose-500" />
          <p>
            O cadastro de
            <span class="font-semibold text-highlighted">{{ nomeDelete }}</span>
            será removido e todas as permissões de zona dessa pessoa serão apagadas.
          </p>
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
          label="Excluir cadastro"
          color="error"
          icon="i-lucide-trash-2"
          :loading="isLoading"
          @click="deleteCadastro" />
      </template>
    </UModal>

    <UModal
      v-model:open="modalLevel"
      title="Alterar nível de acesso"
      description="O nível define o que a pessoa pode fazer no sistema."
      :ui="{ footer: 'justify-end' }">
      <template #body>
        <div class="space-y-4">
          <div class="rounded-lg border border-default bg-muted/40 p-4">
            <p class="text-sm text-toned">
              <span class="font-semibold text-highlighted">{{ nomeLevel }}</span>
              passa de
              <span class="data font-semibold text-highlighted">{{ level }}</span>
              para
              <span class="data font-semibold text-primary">{{ novoLevel }}</span>
              .
            </p>
          </div>

          <UAlert
            v-if="novoLevel === 'admin'"
            color="warning"
            variant="soft"
            icon="i-lucide-shield-alert"
            title="Acesso total"
            description="Administradores enxergam e editam todas as zonas, gerenciam usuários, cadastros e solicitações." />

          <UAlert
            v-else
            color="info"
            variant="soft"
            icon="i-lucide-info"
            title="Acesso por zona"
            description="Membros só enxergam as zonas em que foram adicionados com leitura ou escrita." />
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
              modalLevel = false
            }
          " />
        <UButton
          label="Alterar nível"
          icon="i-lucide-check"
          :loading="isLoading"
          @click="updateLevel" />
      </template>
    </UModal>

    <UModal
      v-model:open="modalEdit"
      title="Editar cadastro"
      :description="`Dados da conta de ${displayName}`"
      :ui="{ footer: 'justify-end', content: 'max-w-lg' }">
      <template #body>
        <UForm
          :schema="cadastroEditSchema"
          :state="editCadastro"
          class="space-y-5"
          @submit="editar">
          <UFormField>
            <UFileUpload v-slot="{ open, removeFile }" v-model="foto" accept="image/*">
              <div class="flex flex-col items-center gap-4">
                <button type="button" class="group relative cursor-pointer" @click="open()">
                  <UAvatar
                    size="3xl"
                    :src="foto ? createObjectUrl(foto) : `${baseUrl}/file/${editCadastro.foto}`"
                    icon="i-lucide-user"
                    class="ring-2 ring-default transition-all group-hover:scale-105 group-hover:ring-primary" />

                  <div
                    class="absolute inset-0 flex items-center justify-center rounded-full bg-black/50 opacity-0 backdrop-blur-sm transition-opacity group-hover:opacity-100">
                    <UIcon name="i-lucide-camera" class="size-6 text-white" />
                  </div>
                </button>

                <div class="flex gap-2">
                  <UButton
                    label="Trocar foto"
                    icon="i-lucide-upload"
                    color="neutral"
                    variant="outline"
                    size="sm"
                    @click="open()" />

                  <UButton
                    v-if="foto"
                    label="Remover"
                    icon="i-lucide-trash-2"
                    color="error"
                    variant="soft"
                    size="sm"
                    @click="removeFile()" />
                </div>
              </div>
            </UFileUpload>
          </UFormField>

          <UFormField label="Nome" name="nome" required>
            <UInput
              v-model="editCadastro.nome"
              icon="i-lucide-user"
              placeholder="Nome completo"
              class="w-full" />
          </UFormField>

          <div class="space-y-3 rounded-lg border border-default bg-muted/40 p-4">
            <p class="text-xs text-dimmed">
              Deixe os campos de senha em branco para manter a senha atual.
            </p>

            <UFormField label="Nova senha" name="senha">
              <UInput
                v-model="editCadastro.senha"
                :color="editCadastro.senha ? color : undefined"
                icon="i-lucide-key"
                :type="show ? 'text' : 'password'"
                placeholder="Nova senha"
                :ui="{ trailing: 'pe-1' }"
                class="w-full">
                <template #trailing>
                  <UButton
                    color="neutral"
                    variant="link"
                    size="sm"
                    :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                    :aria-label="show ? 'Esconder senha' : 'Mostrar senha'"
                    :aria-pressed="show"
                    @click="
                      () => {
                        show = !show
                      }
                    " />
                </template>
              </UInput>
            </UFormField>

            <UFormField label="Confirmar senha" name="confirmar_senha">
              <UInput
                v-model="confirmSenha"
                :color="confirmSenha && confirmSenha !== editCadastro.senha ? 'error' : undefined"
                :type="showConfirm ? 'text' : 'password'"
                placeholder="Confirme a nova senha"
                icon="i-lucide-key"
                :ui="{ trailing: 'pe-1' }"
                class="w-full">
                <template #trailing>
                  <UButton
                    color="neutral"
                    variant="link"
                    size="sm"
                    :icon="showConfirm ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                    :aria-label="showConfirm ? 'Esconder senha' : 'Mostrar senha'"
                    :aria-pressed="showConfirm"
                    @click="
                      () => {
                        showConfirm = !showConfirm
                      }
                    " />
                </template>
              </UInput>
              <p
                v-if="confirmSenha && confirmSenha !== editCadastro.senha"
                class="mt-1 text-xs text-error">
                As senhas não coincidem.
              </p>
            </UFormField>

            <template v-if="editCadastro.senha">
              <UProgress :color="color" :model-value="score" :max="4" size="sm" />

              <p class="text-xs font-medium text-toned">{{ text }}. Deve conter:</p>

              <ul class="space-y-1" aria-label="Requisitos da senha">
                <li
                  v-for="req of strength"
                  :key="req.text"
                  class="flex items-center gap-1.5"
                  :class="req.met ? 'text-success' : 'text-dimmed'">
                  <UIcon
                    :name="req.met ? 'i-lucide-circle-check' : 'i-lucide-circle-x'"
                    class="size-3.5 shrink-0" />
                  <span class="text-xs">{{ req.text }}</span>
                </li>
              </ul>
            </template>
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
              modalEdit = false
            }
          " />
        <UButton
          label="Salvar alterações"
          icon="i-lucide-check"
          :loading="isLoading"
          @click="editar" />
      </template>
    </UModal>
  </UContainer>
</template>
