<script setup lang="ts">
const toast = useToast()

const state = ref<Auth>({ email: '', password: '' })

function login(){
  toast.add({ title: 'Success', description: 'The form has been submitted.', color: 'success' })
}
const show = ref(false)
</script>

<template>
  <section>
    <UCard class="max-w-md mx-auto mt-16 p-6">
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
              <UButton color="neutral" variant="link" size="sm" :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'" :aria-label="show ? 'Esconder senha' : 'Ver senha'" :aria-pressed="show" aria-controls="password" @click="show = !show" />
            </template>
          </UInput>
        </UFormField>

        <UButton type="submit" class="flex justify-center w-full mt-10">
          Login
        </UButton>
      </UForm>
    </UCard>
  </section>
</template>
