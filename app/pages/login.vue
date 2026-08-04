<script setup lang="ts">
  import { safeParse } from 'valibot'

  definePageMeta({ layout: false })
  useHead({ title: 'Login' })

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()
  const { authSchema, cadastroSchema } = useCadastroSchema()

  const state = ref<AuthSchema>({ email: '', senha: '' })

  const show = ref(false)
  const showConfirm = ref(false)
  const modal = ref(false)

  async function login(): Promise<void> {
    start()

    const body = safeParse(authSchema, state.value)
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/login', { method: 'POST', body: body.output }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao entrar',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
    })

    if (!res) return finish({ error: true })

    toast.add({ title: res.message, icon: 'i-lucide-shield-check', color: 'success' })
    await navigateTo('/')
    finish()
  }

  const cadastro = ref<CadastroSchema>({ email: '', senha: '', foto: '', nome: '' })
  const confirmSenha = ref('')
  const foto = ref<File | null>(null)

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
    requisitos.map((req) => ({ met: req.regex.test(cadastro.value.senha), text: req.text })),
  )

  const score = computed(() => strength.value.filter((req) => req.met).length)

  const color = computed(() => {
    if (score.value === 0) return 'neutral'
    if (score.value <= 1) return 'error'
    if (score.value <= 3) return 'warning'
    return 'success'
  })

  const text = computed(() => {
    if (score.value === 0) return 'Digite uma senha'
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

    if (res) cadastro.value.foto = res.message
  })

  async function cadastrar(): Promise<void> {
    start()

    if (cadastro.value.senha !== confirmSenha.value) {
      toast.add({ title: 'As senhas não coincidem', icon: 'i-lucide-shield-alert', color: 'error' })
      return finish({ error: true })
    }

    const body = safeParse(cadastroSchema, cadastro.value)
    if (!body.success) {
      for (const e of body.issues) {
        toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
      }
      return finish({ error: true })
    }

    const res = await $api('/cadastro', { method: 'POST', body: body.output }).catch((error) => {
      toast.add({
        title: error?.data?.message || error?.message || 'Erro ao criar o cadastro',
        icon: 'i-lucide-shield-alert',
        color: 'error',
      })
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
  <div class="relative flex min-h-dvh items-center justify-center overflow-hidden px-4 py-10">
    <div class="pointer-events-none absolute inset-0 console-grid opacity-70" />

    <div class="relative w-full max-w-md">
      <div class="mb-8 flex flex-col items-center gap-3">
        <NuxtImg src="/logo.png" alt="SanchezDNS" width="72" height="72" />
        <p class="text-2xl font-bold tracking-tight text-highlighted">
          Sanchez<span class="text-primary">DNS</span>
        </p>
        <p class="text-center text-sm text-muted">
          Gerenciador de servidores DNS autoritativos com PowerDNS
        </p>
      </div>

      <div class="console-rail overflow-hidden rounded-2xl border border-default bg-default p-7">
        <UForm :schema="authSchema" :state="state" class="space-y-5" @submit="login">
          <div class="space-y-1.5">
            <p class="text-[11px] font-semibold tracking-[0.18em] text-dimmed uppercase">Acesso</p>
            <h1 class="text-xl font-bold text-highlighted">Entrar na sua conta</h1>
          </div>

          <UFormField label="Email" name="email">
            <UInput
              v-model="state.email"
              icon="i-lucide-mail"
              autocomplete="email"
              placeholder="voce@exemplo.com"
              class="w-full"
              :ui="{ base: 'data' }" />
          </UFormField>

          <UFormField label="Senha" name="senha">
            <UInput
              v-model="state.senha"
              class="w-full"
              icon="i-lucide-lock"
              autocomplete="current-password"
              :type="show ? 'text' : 'password'"
              :ui="{ trailing: 'pe-1' }">
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

          <UButton
            type="submit"
            block
            size="lg"
            icon="i-lucide-log-in"
            label="Entrar"
            :loading="isLoading" />
        </UForm>

        <div class="mt-6 border-t border-default pt-5 text-center">
          <p class="text-xs text-muted">
            Ainda não tem conta?
            <UButton
              variant="link"
              size="xs"
              class="px-1"
              label="Cadastre-se"
              @click="
                () => {
                  modal = true
                }
              " />
          </p>
        </div>
      </div>
    </div>

    <UModal
      v-model:open="modal"
      title="Crie seu cadastro"
      description="Sua solicitação passa por aprovação de um administrador antes do primeiro acesso."
      :ui="{ footer: 'justify-end', content: 'max-w-lg' }">
      <template #body>
        <UForm :schema="cadastroSchema" :state="cadastro" class="space-y-5" @submit="cadastrar">
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
                  <p class="text-sm font-medium text-toned">
                    {{ foto ? foto.name : 'Escolha uma foto de perfil' }}
                  </p>
                  <p class="text-xs text-dimmed">JPG, PNG ou GIF</p>
                </div>

                <div class="flex gap-2">
                  <UButton
                    :label="foto ? 'Trocar foto' : 'Selecionar foto'"
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
              v-model="cadastro.nome"
              icon="i-lucide-user"
              placeholder="Seu nome completo"
              class="w-full" />
          </UFormField>

          <UFormField label="Email" name="email" required>
            <UInput
              v-model="cadastro.email"
              icon="i-lucide-mail"
              placeholder="voce@exemplo.com"
              class="w-full"
              :ui="{ base: 'data' }" />
          </UFormField>

          <div class="space-y-3 rounded-lg border border-default bg-muted/40 p-4">
            <UFormField label="Senha" name="senha" required>
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
                :color="confirmSenha && confirmSenha !== cadastro.senha ? 'error' : undefined"
                :type="showConfirm ? 'text' : 'password'"
                placeholder="Confirme sua senha"
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
                v-if="confirmSenha && confirmSenha !== cadastro.senha"
                class="mt-1 text-xs text-error">
                As senhas não coincidem.
              </p>
            </UFormField>

            <UProgress :color="color" :model-value="score" :max="4" size="sm" />

            <p id="password-strength" class="text-xs font-medium text-toned">
              {{ text }}. Deve conter:
            </p>

            <ul class="space-y-1" aria-label="Requisitos da senha">
              <li
                v-for="req of strength"
                :key="req.text"
                class="flex items-center gap-1.5"
                :class="req.met ? 'text-success' : 'text-dimmed'">
                <UIcon
                  :name="req.met ? 'i-lucide-circle-check' : 'i-lucide-circle-x'"
                  class="size-3.5 shrink-0" />
                <span class="text-xs">
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
          color="neutral"
          variant="ghost"
          :loading="isLoading"
          @click="
            () => {
              modal = false
            }
          " />
        <UButton
          label="Criar cadastro"
          icon="i-lucide-check"
          :loading="isLoading"
          @click="cadastrar" />
      </template>
    </UModal>
  </div>
</template>
