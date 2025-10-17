<script setup lang="ts">
definePageMeta({
  layout: false,
})
useHead({ title: 'Login' })

useSeoMeta({ description: 'Login page' })

defineOgImageComponent('Login', { title: 'Login page' })

const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()
const { setUserSession } = useUserSession()

const state = ref<Auth>({ email: '', password: '' })

function checkStrength(str: string){
  const requirements = [
    { regex: /.{8,}/, text: 'At least 8 characters' },
    { regex: /\d/, text: 'At least 1 number' },
    { regex: /[a-z]/, text: 'At least 1 lowercase letter' },
    { regex: /[A-Z]/, text: 'At least 1 uppercase letter' },
  ]

  return requirements.map(req => ({ met: req.regex.test(str), text: req.text }))
}

const strength = computed(() => checkStrength(state.value.password))
const score = computed(() => strength.value.filter(req => req.met).length)

const color = computed(() => {
  if(score.value === 0) return 'neutral'
  if(score.value <= 1) return 'error'
  if(score.value <= 2) return 'warning'
  if(score.value === 3) return 'warning'
  return 'success'
})

const text = computed(() => {
  if(score.value === 0) return 'Enter a password'
  if(score.value <= 2) return 'Weak password'
  if(score.value === 3) return 'Medium password'
  return 'Strong password'
})

const show = ref(false)

const modal = ref(false)

async function login(){
  start()

  const body = AuthSchema.safeParse(state.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string, token: string, isAdmin: boolean }>('/server/login', { method: 'post', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  setUserSession({ username: state.value.email, admin: res.isAdmin, token: res.token })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  await navigateTo('/start')
}

async function register(){
  start()

  const body = AuthSchema.safeParse(state.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/user', { method: 'PUT', body: { ...body.data, level: 'user' } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  modal.value = false
}
</script>

<template>
  <UContainer class="flex items-center justify-center h-screen">
    <div class="flex flex-col items-center gap-6 w-full max-w-md p-6">
      <div class="flex flex-col items-center ">
        <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="92" />
        <p class="mt-2 text-2xl font-bold">
          Sanchez <span class="text-green-500">DNS</span>
        </p>
      </div>
      <UCard class="w-full">
        <template #header>
          <h2 class="text-lg font-medium">
            Login
          </h2>
        </template>

        <UForm :schema="AuthSchema" :state="state" class="space-y-4" @submit="login">
          <UFormField label="Email" name="email">
            <UInput v-model="state.email" icon="i-lucide-mail" class="w-full" />
          </UFormField>

          <UFormField label="Password" name="password">
            <UInput v-model="state.password" class="w-full" icon="i-lucide-lock" :type="show ? 'text' : 'password'">
              <template #trailing>
                <UButton color="neutral" variant="link" size="sm" :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'" :aria-label="show ? 'Hide password' : 'Show password'" :aria-pressed="show" aria-controls="password" @click="show = !show" />
              </template>
            </UInput>
          </UFormField>
          <div class="flex justify-end">
            <UButton variant="ghost" size="sm" @click="modal = true">
              Register...
            </UButton>
          </div>

          <UButton type="submit" class="flex justify-center w-full mt-5">
            Login
          </UButton>
        </UForm>
      </UCard>
    </div>

    <UModal v-model:open="modal" title="Create User" description="Create your account to enjoy the system" :ui="{ footer: 'justify-end' }">
      <template #body>
        <UForm :schema="AuthSchema" :state="state" class="space-y-4" @submit="register">
          <UFormField label="Email" name="email">
            <UInput v-model="state.email" icon="i-lucide-mail" class="w-full" />
          </UFormField>

          <div class="space-y-2">
            <UFormField label="Password">
              <UInput v-model="state.password" placeholder="Password" :color="color" :type="show ? 'text' : 'password'" :aria-invalid="score < 4" aria-describedby="password-strength" :ui="{ trailing: 'pe-1' }" class="w-full">
                <template #trailing>
                  <UButton color="neutral" variant="link" size="sm" :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'" :aria-label="show ? 'Hide password' : 'Show password'" :aria-pressed="show" aria-controls="password" @click="show = !show" />
                </template>
              </UInput>
            </UFormField>
            <UProgress :color="color" :indicator="text" :model-value="score" :max="4" size="sm" />
            <p id="password-strength" class="text-sm font-medium">
              {{ text }}. Must contain:
            </p>
            <ul class="space-y-1" aria-label="Password requirements">
              <li v-for="(req, index) in strength" :key="index" class="flex items-center gap-0.5" :class="req.met ? 'text-success' : 'text-muted'">
                <UIcon :name="req.met ? 'i-lucide-circle-check' : 'i-lucide-circle-x'" class="size-4 shrink-0" />
                <span class="text-xs font-light">
                  {{ req.text }}
                  <span class="sr-only">
                    {{ req.met ? ' - Requirement met' : ' - Requirement not met' }}
                  </span>
                </span>
              </li>
            </ul>
          </div>
        </UForm>
      </template>

      <template #footer>
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
        <UButton label="Confirm" :loading="isLoading" @click="register" />
      </template>
    </UModal>
  </UContainer>
</template>
