<script setup lang="ts">
useHead({ title: 'Configuration server' })

useSeoMeta({ description: 'Configuration server page' })

defineOgImageComponent('Techs', { title: 'Configuration server page' })

const { optionSelected, refreshConnections } = await useConnection()
const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()

const { data, refresh } = await useFetch<{ connection: { _id: string, name: string, host: string, serverId: string, createdAt: string, updatedAt: string }, server: { id: string, daemonType: string, backend: string, databasePath: string, listeningAddress: string, webserver: string, api: boolean } }>('/server/api/connection', { method: 'GET', query: { connection: optionSelected.value } })

const isEditing = ref(false)

const stateConnection = ref<EditConnectionType>({ name: data.value?.connection.name || '', host: data.value?.connection.host || '', serverId: data.value?.server.id || '' })

function startEditing(){
  if(data.value){
    stateConnection.value = { name: data.value.connection.name, host: data.value.connection.host, serverId: data.value.server.id }
    isEditing.value = true
  }
}

function cancelEditing(){
  isEditing.value = false
  if(data.value) stateConnection.value = { name: data.value.connection.name, host: data.value.connection.host, serverId: data.value.server.id }
}

const modal = ref(false)

function confirmEditing(){
  if((stateConnection.value.host !== data.value?.connection.host) || (stateConnection.value.serverId !== data.value?.server.id)){
    return modal.value = true
  }
  editConnection()
}

const confirmHost = ref('')

async function editConnection(){
  start()

  const body = EditConnectionSchema.safeParse(stateConnection.value)

  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/connection', { method: 'PATCH', body: body.data, query: { connection: optionSelected.value } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  confirmHost.value = ''
  await refreshConnections()
  await refresh()
  isEditing.value = false
}

const modalDelete = ref(false)

async function deleteConnection(){
  start()

  if(data.value && (confirmHost.value !== data.value.connection.host)){
    toast.add({ title: 'The host does not match', icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/connection', { method: 'DELETE', query: { connection: optionSelected.value } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  await refreshConnections()
  optionSelected.value = ''
  await navigateTo('/start')
}

watch(modal, newVal => {
  if(!newVal) confirmHost.value = ''
})
watch(modalDelete, newVal => {
  if(!newVal) confirmHost.value = ''
})
</script>

<template>
  <UContainer>
    <div v-if="data" class="py-10 space-y-8">
      <h1 class="text-3xl font-bold text-gray-100">
        Server Configuration
      </h1>

      <UCard>
        <template #header>
          <h2 class="text-gray-200 font-medium">
            Server Information
          </h2>
        </template>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div class="rounded-lg border border-[var(--ui-border)] p-4 bg-[var(--ui-bg-elevated)]">
            <div class="text-xs text-[var(--ui-text-muted)]">
              Server ID
            </div>
            <div class="mt-1 font-medium text-[var(--ui-text)]">
              {{ data.server.id }}
            </div>
          </div>
          <div class="rounded-lg border border-[var(--ui-border)] p-4 bg-[var(--ui-bg-elevated)]">
            <div class="text-xs text-[var(--ui-text-muted)]">
              Daemon Type
            </div>
            <div class="mt-1 font-medium text-[var(--ui-text)]">
              {{ data.server.daemonType }}
            </div>
          </div>
          <div class="rounded-lg border border-[var(--ui-border)] p-4 bg-[var(--ui-bg-elevated)]">
            <div class="text-xs text-[var(--ui-text-muted)]">
              Backend
            </div>
            <div class="mt-1 font-medium text-[var(--ui-text)]">
              {{ data.server.backend }}
            </div>
          </div>
          <div class="rounded-lg border border-[var(--ui-border)] p-4 bg-[var(--ui-bg-elevated)]">
            <div class="text-xs text-[var(--ui-text-muted)]">
              Database Path
            </div>
            <div class="mt-1 font-medium text-[var(--ui-text)] break-all">
              {{ data.server.databasePath }}
            </div>
          </div>
          <div class="rounded-lg border border-[var(--ui-border)] p-4 bg-[var(--ui-bg-elevated)]">
            <div class="text-xs text-[var(--ui-text-muted)]">
              Listening Address
            </div>
            <div class="mt-1 font-medium text-[var(--ui-text)]">
              {{ data.server.listeningAddress }}
            </div>
          </div>
          <div class="rounded-lg border border-[var(--ui-border)] p-4 bg-[var(--ui-bg-elevated)]">
            <div class="text-xs text-[var(--ui-text-muted)]">
              Webserver
            </div>
            <div class="mt-1 font-medium text-[var(--ui-text)]">
              {{ data.server.webserver }}
            </div>
          </div>
          <div class="rounded-lg border border-[var(--ui-border)] p-4 bg-[var(--ui-bg-elevated)]">
            <div class="text-xs text-[var(--ui-text-muted)]">
              API
            </div>
            <div class="mt-1">
              <UBadge variant="soft">
                {{ data.server.api ? 'Enabled' : 'Disabled' }}
              </UBadge>
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h2 class="text-gray-200 font-medium">
              Connection Information
            </h2>
            <div class="flex items-center gap-2">
              <UButton v-if="isEditing" icon="i-heroicons-x-mark" color="error" size="xs" :loading="isLoading" label="Cancel" @click="cancelEditing" />
              <UButton v-if="isEditing" icon="i-heroicons-check" size="xs" :loading="isLoading" label="Save" @click="confirmEditing()" />
            </div>
          </div>
        </template>
        <UForm :schema="EditConnectionSchema" :state="stateConnection" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mx-auto" :disabled="!isEditing">
          <UFormField label="Connection Name">
            <UInput v-model="stateConnection.name" :disabled="!isEditing" placeholder="Name" />
          </UFormField>

          <UFormField label="Server ID">
            <UInput v-model="stateConnection.serverId" :disabled="!isEditing" placeholder="localhost" />
          </UFormField>

          <UFormField label="Host">
            <UInput v-model="stateConnection.host" :disabled="!isEditing" placeholder="152.84.120.200" />
          </UFormField>
        </UForm>

        <template #footer>
          <div class="flex items-center gap-10">
            <div>
              <p class="text-gray-200 font-medium">
                Created at
              </p>
              <NuxtTime :datetime="data.connection.createdAt" day="2-digit" month="2-digit" year="2-digit" hour="2-digit" minute="2-digit" class="text-sm text-gray-500" />
            </div>
            <div>
              <p class="text-gray-200 font-medium">
                Last updated
              </p>
              <NuxtTime :datetime="data.connection.updatedAt" day="2-digit" month="2-digit" year="2-digit" hour="2-digit" minute="2-digit" class="text-sm text-gray-500" />
            </div>
          </div>
        </template>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="text-gray-200 font-medium">
            Server Actions
          </h2>
        </template>

        <div class="flex flex-wrap gap-3">
          <UButton :loading="isLoading" icon="i-lucide-bar-chart-3" color="neutral" to="/srv/statistics" label="View Statistics" />
          <UButton :loading="isLoading" icon="i-lucide-globe" to="/zones/dashboard" label="View Zones" color="info" />
          <UButton :loading="isLoading" icon="i-lucide-pen" label="Edit Connection" :disabled="isEditing" @click="startEditing" />
          <UButton :loading="isLoading" icon="i-lucide-trash" label="Delete Connection" color="error" @click="modalDelete = true" />
        </div>
      </UCard>
    </div>

    <div v-else class="flex items-center justify-center h-64">
      <UProgress animation="carousel" />
    </div>

    <UModal v-model:open="modal" title="Danger" description="You have changed critical connection data, are you sure you want to continue?" :ui="{ footer: 'justify-end' }">
      <template #body>
        <div class="flex items-center">
          <p class="text-gray-200">
            If you are sure you want to continue, write the host of your connection below <span class="font-bold">'{{ stateConnection.host }}'</span>.
          </p>
          <UTooltip text="The host and the server ID are critical data, if you change them you may not be able to connect to your server." :delay-duration="0" arrow>
            <UButton icon="i-lucide-info" color="neutral" variant="subtle" class="rounded-full" />
          </UTooltip>
        </div>
        <UInput v-model="confirmHost" class="mt-2 w-full" color="error" placeholder="Host" />
      </template>

      <template #footer>
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
        <UButton label="Confirm" color="error" :loading="isLoading" :disabled="confirmHost !== stateConnection.host" @click="editConnection" />
      </template>
    </UModal>

    <UModal v-model:open="modalDelete" title="Danger" description="You are about to delete a connection, this action cannot be undone." :ui="{ footer: 'justify-end' }">
      <template #body>
        <p class="text-gray-200">
          If you are sure you want to continue, write the host of your connection below <span class="font-bold">'{{ stateConnection.host }}'</span>.
        </p>
        <UInput v-model="confirmHost" class="mt-2 w-full" color="error" placeholder="Host" />
      </template>

      <template #footer>
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modalDelete = false" />
        <UButton label="Confirm" color="error" :loading="isLoading" :disabled="confirmHost !== stateConnection.host" @click="deleteConnection" />
      </template>
    </UModal>
  </UContainer>
</template>
