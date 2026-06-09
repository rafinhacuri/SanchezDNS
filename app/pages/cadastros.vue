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

  const { data, refresh } = await useFetch<CadastroResponse>('/server/api/cadastros', {
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

    const body = safeParse(IdSchema, {
      id: idDelete.value ?? '',
    })
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/cadastro', {
      method: 'DELETE',
      body: body.output,
    }).catch((error) => {
      toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
    })

    if (!res) return finish({ error: true })

    refresh()
    toast.add({ title: res.message, icon: 'i-lucide-shield-check', color: 'success' })
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
  const modalLevel = ref(false)
  const level = ref('')

  function openLevel(id: string, nome: string, lvl: string): void {
    idLevel.value = id
    nomeLevel.value = nome
    level.value = lvl

    modalLevel.value = true
  }

  async function updateLevel(): Promise<void> {
    start()

    const body = safeParse(IdSchema, {
      id: idLevel.value ?? '',
    })
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/level', {
      method: 'PATCH',
      body: body.output,
    }).catch((error) => {
      toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
    })

    if (!res) return finish({ error: true })

    refresh()
    toast.add({ title: res.message, icon: 'i-lucide-shield-check', color: 'success' })
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

  function checkStrength(str: string): { met: boolean; text: string }[] {
    const requirements = [
      // oxlint-disable-next-line require-unicode-regexp
      { regex: /.{8,}/, text: 'Pelo menos 8 caracteres' },
      // oxlint-disable-next-line require-unicode-regexp
      { regex: /\d/, text: 'Pelo menos 1 número' },
      // oxlint-disable-next-line require-unicode-regexp
      { regex: /[a-z]/, text: 'Pelo menos 1 letra minúscula' },
      // oxlint-disable-next-line require-unicode-regexp
      { regex: /[A-Z]/, text: 'Pelo menos 1 letra maiúscula' },
    ]

    return requirements.map((req) => ({ met: req.regex.test(str), text: req.text }))
  }

  const modalEdit = ref(false)

  const displayName = ref('')
  const editCadastro = ref({
    id: '',
    senha: '',
    foto: '',
    nome: '',
  })

  function openEdit(cadastro: Cadastro): void {
    editCadastro.value = { ...cadastro }
    displayName.value = cadastro.nome

    modalEdit.value = true
  }

  const confirmSenha = ref('')

  const strength = computed(() => checkStrength(editCadastro.value?.senha || ''))
  const score = computed(() => strength.value.filter((req) => req.met).length)

  const color = computed(() => {
    if (score.value === 0) return 'neutral'
    if (score.value <= 1) return 'error'
    if (score.value <= 2) return 'warning'
    if (score.value === 3) return 'warning'
    return 'success'
  })

  const text = computed(() => {
    if (score.value === 0) return 'Digite uma senha'
    if (score.value <= 2) return 'Senha fraca'
    if (score.value === 3) return 'Senha média'
    return 'Senha forte'
  })

  const show = ref(false)
  const showConfirm = ref(false)

  const foto = ref<File | null>(null)

  function createObjectUrl(file: File): string {
    return URL.createObjectURL(file)
  }

  watch(foto, async (file) => {
    if (file) {
      const formData = new FormData()
      formData.append('file', file)

      const res = await $fetch<GoRes>('/server/api/file', {
        method: 'PUT',
        body: formData,
      }).catch((error) => {
        toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
      })

      if (res) editCadastro.value.foto = res.message
    }
  })

  async function editar(): Promise<void> {
    start()

    if (editCadastro.value.senha !== confirmSenha.value) {
      toast.add({ title: 'As senhas não coincidem', icon: 'i-lucide-shield-alert', color: 'error' })
      return finish({ error: true })
    }

    const body = safeParse(CadastroEditSchema, editCadastro.value)
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/cadastro', {
      method: 'PUT',
      body: body.output,
    }).catch((error) => {
      toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-shield-check', color: 'success' })
    refresh()
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
              @click="openEdit(cadastro)"
              block />
            <UButton
              :loading="isLoading"
              label="Alterar Level"
              color="warning"
              icon="i-lucide-user-pen"
              block
              @click="openLevel(cadastro.id, cadastro.nome, cadastro.level)"
              variant="soft" />
            <UButton
              v-if="cadastro.level !== 'admin'"
              :loading="isLoading"
              label="Excluir"
              color="error"
              icon="i-lucide-trash-2"
              variant="soft"
              @click="openDelete(cadastro.id, cadastro.nome)"
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

    <UModal v-model:open="modalDelete" :title="`Excluir cadastro de ${nomeDelete}`">
      <template #body>
        <p class="text-sm text-gray-500 dark:text-gray-400">
          Tem certeza que deseja excluir este cadastro? Esta ação é irreversível.
        </p>
      </template>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <UButton
            label="Cancelar"
            variant="outline"
            color="neutral"
            @click="modalDelete = false" />
          <UButton :loading="isLoading" label="Excluir" color="error" @click="deleteCadastro" />
        </div>
      </template>
    </UModal>

    <UModal v-model:open="modalLevel" :title="`Alterar level de ${nomeLevel}`">
      <template #body>
        <p class="text-sm text-gray-500 dark:text-gray-400">
          Tem certeza que deseja alterar o level deste cadastro para
          {{ level == 'admin' ? 'Usuário' : 'Administrador' }}?
        </p>
      </template>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <UButton label="Cancelar" variant="outline" color="neutral" @click="modalLevel = false" />
          <UButton
            :loading="isLoading"
            label="confirmar"
            color="warning"
            variant="soft"
            @click="updateLevel" />
        </div>
      </template>
    </UModal>

    <UModal
      v-model:open="modalEdit"
      :title="`Editar cadastro de ${displayName}`"
      description="Faça as alterações desejadas no cadastro. Certifique-se de fornecer informações válidas."
      :ui="{ footer: 'justify-end' }">
      <template #body>
        <UForm :schema="CadastroSchema" :state="editCadastro" class="space-y-4">
          <UFormField>
            <UFileUpload v-slot="{ open, removeFile }" v-model="foto" accept="image/*">
              <div class="flex flex-col items-center gap-4">
                <button type="button" class="group relative cursor-pointer" @click="open()">
                  <UAvatar
                    size="3xl"
                    :src="foto ? createObjectUrl(foto) : undefined"
                    icon="i-lucide-user"
                    class="ring-2 ring-default transition-all group-hover:scale-105 group-hover:ring-primary" />

                  <div
                    class="absolute inset-0 flex items-center justify-center rounded-full bg-black/50 opacity-0 backdrop-blur-sm transition-opacity group-hover:opacity-100">
                    <UIcon name="i-lucide-camera" class="size-6 text-white" />
                  </div>
                </button>

                <div class="text-center">
                  <p class="text-sm font-medium">
                    {{ foto ? foto.name : 'Escolha uma foto de perfil' }}
                  </p>

                  <p class="text-xs text-muted">JPG, PNG ou GIF</p>
                </div>

                <div class="flex gap-2">
                  <UButton
                    :label="foto ? 'Trocar foto' : 'Selecionar foto'"
                    icon="i-lucide-upload"
                    color="neutral"
                    variant="outline"
                    @click="open()" />

                  <UButton
                    v-if="foto"
                    label="Remover"
                    icon="i-lucide-trash-2"
                    color="error"
                    variant="soft"
                    @click="removeFile()" />
                </div>
              </div>
            </UFileUpload>
          </UFormField>

          <UFormField label="Nome" name="nome">
            <UInput
              v-model="editCadastro.nome"
              icon="i-lucide-user"
              placeholder="Digite seu nome completo"
              class="w-full" />
          </UFormField>

          <div class="space-y-2">
            <UFormField label="Senha" name="senha">
              <UInput
                v-model="editCadastro.senha"
                :color="color"
                icon="i-lucide-key"
                :type="show ? 'text' : 'password'"
                :aria-invalid="score < 4"
                placeholder="Digite uma senha forte"
                aria-describedby="password-strength"
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
                    aria-controls="password"
                    @click="show = !show" />
                </template>
              </UInput>
            </UFormField>
            <UFormField label="Confirmar Senha" name="confirmar_senha">
              <UInput
                v-model="confirmSenha"
                :color="confirmSenha && confirmSenha !== editCadastro.senha ? 'error' : undefined"
                :type="showConfirm ? 'text' : 'password'"
                placeholder="Confirme sua senha"
                icon="i-lucide-key"
                :aria-invalid="confirmSenha && confirmSenha !== editCadastro.senha"
                aria-describedby="confirm-password-error"
                class="w-full">
                <template #trailing>
                  <UButton
                    color="neutral"
                    variant="link"
                    size="sm"
                    :icon="showConfirm ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                    :aria-label="showConfirm ? 'Esconder senha' : 'Mostrar senha'"
                    :aria-pressed="showConfirm"
                    aria-controls="password"
                    @click="showConfirm = !showConfirm" />
                </template>
              </UInput>
              <p
                v-if="confirmSenha && confirmSenha !== editCadastro.senha"
                id="confirm-password-error"
                class="text-sm text-error">
                As senhas não coincidem.
              </p>
            </UFormField>
            <UProgress :color="color" :indicator="text" :model-value="score" :max="4" size="sm" />
            <p id="password-strength" class="text-sm font-medium">{{ text }}. Deve conter:</p>
            <ul class="space-y-1" aria-label="Requisitos da senha">
              <li
                v-for="(req, index) in strength"
                :key="index"
                class="flex items-center gap-0.5"
                :class="req.met ? 'text-success' : 'text-muted'">
                <UIcon
                  :name="req.met ? 'i-lucide-circle-check' : 'i-lucide-circle-x'"
                  class="size-4 shrink-0" />
                <span class="text-xs font-light">
                  {{ req.text }}
                  <span class="sr-only">
                    {{ req.met ? ' - Requisito cumprido' : ' - Requisito não cumprido' }}
                  </span>
                </span>
              </li>
            </ul>
          </div>
        </UForm>
      </template>

      <template #footer>
        <UButton
          label="Cancelar"
          :loading="isLoading"
          variant="outline"
          @click="modalEdit = false" />
        <UButton label="Confirmar" :loading="isLoading" @click="editar" />
      </template>
    </UModal>
  </UContainer>
</template>
