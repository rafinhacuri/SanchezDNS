<script setup lang="ts">
  import { safeParse } from 'valibot'

  definePageMeta({
    layout: false,
  })
  useHead({ title: 'Login' })

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()

  const state = ref<Auth>({ email: '', senha: '' })

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

  const cadastro = ref<CadastroType>({
    email: '',
    senha: '',
    foto: '',
    nome: '',
  })

  const confirmSenha = ref('')

  const strength = computed(() => checkStrength(cadastro.value.senha))
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

  const modal = ref(false)

  async function login(): Promise<void> {
    start()

    const body = safeParse(AuthSchema, state.value)
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/login', {
      method: 'post',
      body: body.output,
    }).catch((error) => {
      toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-shield-check', color: 'success' })
    await navigateTo('/')
    finish()
  }

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

      if (res) cadastro.value.foto = res.message
    }
  })

  async function cadastrar(): Promise<void> {
    start()

    if (cadastro.value.senha !== confirmSenha.value) {
      toast.add({ title: 'As senhas não coincidem', icon: 'i-lucide-shield-alert', color: 'error' })
      return finish({ error: true })
    }

    const body = safeParse(CadastroSchema, cadastro.value)
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $fetch<GoRes>('/server/api/cadastro', {
      method: 'post',
      body: body.output,
    }).catch((error) => {
      toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-shield-check', color: 'success' })
    modal.value = false
    finish()
  }

  watch(modal, (open) => {
    if (!open) {
      cadastro.value = { email: '', senha: '', foto: '', nome: '' }
      foto.value = null
      confirmSenha.value = ''
    }
  })
</script>

<template>
  <UContainer class="flex h-screen items-center justify-center">
    <div class="flex w-full max-w-md flex-col items-center gap-6 p-6">
      <div class="flex flex-col items-center">
        <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="92" />
        <p class="mt-2 text-2xl font-bold">Sanchez <span class="text-green-500">DNS</span></p>
      </div>
      <UCard class="w-full">
        <template #header>
          <h2 class="text-lg font-medium">Login</h2>
        </template>

        <UForm :schema="AuthSchema" :state="state" class="space-y-4" @submit="login">
          <UFormField label="Email" name="email">
            <UInput v-model="state.email" icon="i-lucide-mail" class="w-full" />
          </UFormField>

          <UFormField label="Senha" name="senha">
            <UInput
              v-model="state.senha"
              class="w-full"
              icon="i-lucide-lock"
              :type="show ? 'text' : 'password'">
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
          <div class="flex justify-end">
            <UButton variant="ghost" size="sm" @click="modal = true"> Cadastre-se... </UButton>
          </div>

          <UButton type="submit" class="mt-5 flex w-full justify-center"> Login </UButton>
        </UForm>
      </UCard>
    </div>

    <UModal
      v-model:open="modal"
      title="Crie seu cadastro"
      description="Crie sua conta para aproveitar o sistema de monitoramento e gerenciamento de DNS"
      :ui="{ footer: 'justify-end' }">
      <template #body>
        <UForm :schema="CadastroSchema" :state="cadastro" class="space-y-4" @submit="cadastrar">
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
              v-model="cadastro.nome"
              icon="i-lucide-user"
              placeholder="Digite seu nome completo"
              class="w-full" />
          </UFormField>

          <UFormField label="Email" name="email">
            <UInput
              v-model="cadastro.email"
              icon="i-lucide-mail"
              placeholder="Digite seu email"
              class="w-full" />
          </UFormField>

          <div class="space-y-2">
            <UFormField label="Senha" name="senha">
              <UInput
                v-model="cadastro.senha"
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
                :color="confirmSenha && confirmSenha !== cadastro.senha ? 'error' : undefined"
                :type="showConfirm ? 'text' : 'password'"
                placeholder="Confirme sua senha"
                icon="i-lucide-key"
                :aria-invalid="confirmSenha && confirmSenha !== cadastro.senha"
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
                v-if="confirmSenha && confirmSenha !== cadastro.senha"
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
        <UButton label="Cancelar" :loading="isLoading" variant="outline" @click="modal = false" />
        <UButton label="Confirmar" :loading="isLoading" @click="cadastrar" />
      </template>
    </UModal>
  </UContainer>
</template>
