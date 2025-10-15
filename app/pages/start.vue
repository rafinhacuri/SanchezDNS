<script setup lang="ts">
const { user } = useUserSession()

const modal = ref(false)

const state = ref<ConnectionType>({ name: '', host: '', apiKey: '', serverId: '', users: [] })

const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()
const { refreshConnections } = await useConnection()

async function createConnection(){
  start()

  const body = ConnectionSchema.safeParse(state.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/connections', { method: 'post', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  await refreshConnections()
  state.value = { name: '', host: '', apiKey: '', serverId: '', users: [] }
  modal.value = false
}
</script>

<template>
  <UContainer class="flex flex-col items-center justify-center h-full text-center p-4 space-y-4">
    <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="92" class="mt-16" />
    <h1 class="text-2xl font-bold">
      Welcome to SanchezDNS
    </h1>
    <p class="text-lg">
      Please select a DNS connection to get started <strong v-if="user.admin">or create a new one</strong>.
    </p>
    <UButton v-if="user.admin" color="primary" variant="outline" @click="modal = true">
      Create Connection
    </UButton>

    <UModal v-model:open="modal" title="Create Connection" description="Create a new connection to an authoritative server" :ui="{ footer: 'justify-end' }">
      <template #body>
        <UForm :schema="ConnectionSchema" :state="state" class="space-y-4">
          <UFormField label="Name" name="name">
            <UInput v-model="state.name" icon="i-lucide-computer" class="w-full" />
          </UFormField>
          <UFormField label="Host" name="host">
            <UInput v-model="state.host" icon="i-lucide-server" class="w-full" />
          </UFormField>
          <UFormField label="API Key" name="apiKey">
            <UInput v-model="state.apiKey" icon="i-lucide-key" class="w-full" />
          </UFormField>
          <UFormField label="Server ID" name="serverId">
            <UInput v-model="state.serverId" icon="i-lucide-hash" class="w-full" />
          </UFormField>
        </UForm>
      </template>

      <template #footer>
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
        <UButton label="Confirm" :loading="isLoading" @click="createConnection" />
      </template>
    </UModal>
  </UContainer>
</template>
