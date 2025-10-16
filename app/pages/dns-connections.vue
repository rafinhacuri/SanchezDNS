<script setup lang="ts">
useHead({ title: 'DNS Machines' })

useSeoMeta({ description: 'DNS Machines page' })

defineOgImageComponent('DNS Machines', { title: 'DNS Machines page' })

interface FullConnection{
  _id: string,
  name: string,
  host: string,
  apiKey: string,
  serverId: string,
  users: string[],
  createdAt: string,
  updatedAt: string,
}

const { data, refresh } = await useFetch<FullConnection[]>('/server/api/full-connections', { method: 'GET' })

const { optionSelected, refreshConnections } = await useConnection()
const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()

async function openDetails(id: string){
  optionSelected.value = id
  await navigateTo('/srv/statistics')
}

async function openConnection(id: string){
  optionSelected.value = id
  await navigateTo('/zones')
}

async function editConnection(id: string){
  optionSelected.value = id
  await navigateTo('/srv/configuration')
}

const modalDelete = ref(false)
const idToDelete = ref('')
const hostToDelete = computed(() => {
  const conn = data.value?.find(c => c._id === idToDelete.value)
  return conn ? conn.host : ''
})
const confirmHost = ref('')

function openDelete(id: string){
  idToDelete.value = id
  modalDelete.value = true
}

async function deleteConnection(){
  start()

  if(data.value && (confirmHost.value !== hostToDelete.value)){
    toast.add({ title: 'The host does not match', icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/connection', { method: 'DELETE', query: { connection: idToDelete.value } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  await refreshConnections()
  if(optionSelected.value === idToDelete.value) optionSelected.value = ''
  await refresh()
  modalDelete.value = false
}

watch(modalDelete, newVal => {
  if(!newVal){
    idToDelete.value = ''
    confirmHost.value = ''
  }
})

const modal = ref(false)

const state = ref<ConnectionType>({ name: '', host: '', apiKey: '', serverId: '', users: [] })

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
  await refresh()
  modal.value = false
}

watch(modal, newVal => {
  if(!newVal) state.value = { name: '', host: '', apiKey: '', serverId: '', users: [] }
})
</script>

<template>
  <UContainer class="p-4">
    <UButton :loading="isLoading" label="Crie uma nova conexão" variant="outline" class="mb-4" icon="i-lucide-radio-tower" @click="modal = true" />

    <div v-if="data">
      <div v-for="c of data" :key="c._id" class="mb-8">
        <button class="bg-gradient-to-br from-gray-800 to-gray-900 rounded-xl p-6 shadow-lg hover:shadow-2xl transition-all duration-300 hover:scale-[1.02] w-full text-left" @click="openConnection(c._id)">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="text-lg font-semibold text-white flex items-center gap-2">
                <UIcon name="i-heroicons-server" class="text-blue-400" />
                {{ c.name }}
              </h2>
              <p class="text-sm text-gray-400">
                {{ c.host }}
              </p>
            </div>
            <div class="flex items-center gap-2">
              <UButton :loading="isLoading" label="Editar" size="sm" color="info" variant="soft" @click.stop="editConnection(c._id)" />
              <UButton :loading="isLoading" label="Excluir" size="sm" color="error" variant="outline" @click.stop="openDelete(c._id)" />
              <UButton :loading="isLoading" label="Ver detalhes" size="sm" variant="soft" @click.stop="openDetails(c._id)" />
            </div>
          </div>

          <div class="border-t border-gray-700 pt-4">
            <p class="text-sm text-gray-500 mb-2">
              Usuários vinculados:
            </p>
            <UAvatarGroup :max="6">
              <UAvatar v-for="user in c.users" :key="user" :alt="user" :style="{ backgroundColor: `hsl(${Math.random() * 360}, 70%, 55%)` }" />
            </UAvatarGroup>
          </div>
        </button>
      </div>
    </div>

    <UModal v-model:open="modalDelete" title="Danger" description="You are about to delete a connection, this action cannot be undone." :ui="{ footer: 'justify-end' }">
      <template #body>
        <p class="text-gray-200">
          If you are sure you want to continue, write the host of your connection below <span class="font-bold">'{{ hostToDelete }}'</span>.
        </p>
        <UInput v-model="confirmHost" class="mt-2 w-full" color="error" placeholder="Host" />
      </template>

      <template #footer>
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modalDelete = false" />
        <UButton label="Confirm" color="error" :loading="isLoading" :disabled="confirmHost !== hostToDelete" @click="deleteConnection" />
      </template>
    </UModal>

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
          <UFormField label="Server ID (Default: 'localhost')" name="serverId">
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
