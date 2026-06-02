<script setup lang="ts">
  import { safeParse } from 'valibot'

  definePageMeta({
    layout: false,
  })
  useHead({ title: 'Login' })

  useSeoMeta({ description: 'Login page' })

  defineOgImageComponent('Login', { title: 'Login page' })

  const toast = useToast()
  const { isLoading, start, finish } = useLoadingIndicator()

  const state = ref<Auth>({ email: '', senha: '' })

  function checkStrength(str: string): { met: boolean; text: string }[] {
    const requirements = [
      { regex: /.{8,}/, text: 'Pelo menos 8 caracteres' },
      { regex: /\d/, text: 'Pelo menos 1 número' },
      { regex: /[a-z]/, text: 'Pelo menos 1 letra minúscula' },
      { regex: /[A-Z]/, text: 'Pelo menos 1 letra maiúscula' },
    ]

    return requirements.map((req) => ({ met: req.regex.test(str), text: req.text }))
  }

  const strength = computed(() => checkStrength(state.value.senha))
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
        <UForm :schema="AuthSchema" :state="state" class="space-y-4">
          <UFormField label="Email" name="email">
            <UInput v-model="state.email" icon="i-lucide-mail" class="w-full" />
          </UFormField>

          <div class="space-y-2">
            <UFormField label="Senha" name="senha">
              <UInput
                v-model="state.senha"
                placeholder="Senha"
                :color="color"
                :type="show ? 'text' : 'password'"
                :aria-invalid="score < 4"
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
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
        <UButton label="Confirm" :loading="isLoading" />
      </template>
    </UModal>
  </UContainer>
</template>
