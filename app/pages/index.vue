<script setup lang="ts">
definePageMeta({
  layout: false,
})
useHead({ title: 'Login' })

useSeoMeta({ description: 'Login page' })

defineOgImageComponent('Techs', { title: 'Login page' })

const toast = useToast()
const { isLoading } = useLoadingIndicator()

const state = ref<Auth>({ email: '', password: '' })

function login(){
  toast.add({ title: 'Success', description: 'The form has been submitted.', color: 'success' })
}
const show = ref(false)

const modal = ref(false)
</script>

<template>
  <section class="flex items-center justify-center h-screen">
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
        </UForm>
      </template>

      <template #footer>
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
        <UButton label="Confirm" :loading="isLoading" />
      </template>
    </UModal>
  </section>
</template>
