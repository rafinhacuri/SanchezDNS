<script setup lang="ts">
useHead({ title: 'Groups' })

useSeoMeta({ description: 'Groups page' })

defineOgImageComponent('Techs', { title: 'Groups page' })

const toast = useToast()
const { isLoading, start, finish } = useLoadingIndicator()
const { nameServer, usersServer, optionSelected, refreshConnections } = await useConnection()

const { data } = await useFetch<{ _id: string, email: string }[]>('/server/api/users', { method: 'GET' })

const modal = ref(false)
const selectedUser = ref('')

async function addUser(){
  start()
  if(!selectedUser.value){
    toast.add({ title: 'Please select a user', icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/connection/user', { method: 'POST', body: { email: selectedUser.value, connection: optionSelected.value } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  toast.add({ title: res.message, icon: 'i-lucide-check-circle', color: 'success' })
  await refreshConnections()
  modal.value = false
  selectedUser.value = ''
  finish()
}
</script>

<template>
  <section class="max-w-7xl mx-auto p-4 w-full">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">
        {{ nameServer }} Member
      </h1>

      <UButton icon="i-lucide-user-plus" class="text-white" variant="soft" color="primary" @click="modal = true">
        Add Member
      </UButton>
    </div>

    <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-6 shadow-sm">
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100 mb-4">
        Members
      </h2>
      <ul class="divide-y divide-gray-200 dark:divide-gray-700">
        <li v-for="member in usersServer" :key="member" class="flex items-center justify-between py-3">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 flex items-center justify-center rounded-full text-white font-medium shadow-md" :style="{ backgroundColor: `hsl(${Math.random() * 360}, 70%, 55%)` }">
              {{ member.charAt(0).toUpperCase() }}
            </div>
            <span class="text-gray-700 dark:text-gray-200">{{ member }}</span>
          </div>
          <UButton icon="i-lucide-trash" variant="outline" color="error">
            Remove Member
          </UButton>
        </li>
      </ul>
    </div>

    <UModal v-model:open="modal" title="Add User" description="Add a new user to the group" :ui="{ footer: 'justify-end' }">
      <template #body>
        <USelectMenu v-model="selectedUser" :items="data" label-key="email" value-key="email" class="w-full" placeholder="Select a user" />
      </template>

      <template #footer>
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modal = false" />
        <UButton label="Confirm" :loading="isLoading" @click="addUser" />
      </template>
    </UModal>
  </section>
</template>
