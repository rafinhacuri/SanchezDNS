<script setup lang="ts">
import type { DropdownMenuItem, NavigationMenuItem } from '@nuxt/ui'

const route = useRoute()

const { user, clearUserSession } = useUserSession()

const { isLoading, start, finish } = useLoadingIndicator()

const { optionSelected, optionsConnection } = useConnection()

const toast = useToast()

const items = computed<NavigationMenuItem[]>(() => [
  {
    label: 'Zones',
    icon: 'i-lucide-database',
    to: '/zones/dashboard',
    active: route.path === '/zones/dashboard' || route.path.startsWith('/zones/templates'),
    children: [
      {
        label: 'Dashboard',
        icon: 'i-lucide-chart-pie',
        description: 'View and manage all DNS zones.',
        to: '/zones/dashboard',
      },
      {
        label: 'Templates',
        icon: 'i-lucide-file-text',
        description: 'Manage DNS templates.',
        to: '/zones/templates',
      },
    ],
  },
  {
    label: 'Server',
    icon: 'i-lucide-server',
    active: route.path.startsWith('/srv/statistics') || route.path.startsWith('/srv/configuration'),
    children: [
      {
        label: 'Statistics',
        icon: 'i-lucide-bar-chart-3',
        description: 'View server statistics.',
        to: '/srv/statistics',
      },
      {
        label: 'Configuration',
        icon: 'i-lucide-settings',
        description: 'Manage server configuration.',
        to: '/srv/configuration',
      },
    ],
  },
  {
    label: 'Configuration',
    icon: 'i-lucide-settings',
    active: route.path.startsWith('/config/dns-machines') || route.path.startsWith('/config/api-keys'),
    children: [
      {
        label: 'DNS Machines',
        icon: 'i-lucide-computer',
        description: 'Manage DNS machines connections.',
        to: '/config/dns-machines',
      },
      {
        label: 'API Keys',
        icon: 'i-lucide-key',
        description: 'Manage API keys.',
        to: '/config/api-keys',
      },
    ],
  },
  ...(user.value?.admin ? [{ label: 'Users', icon: 'i-lucide-users', to: '/users', active: route.path.startsWith('/users') }] : []),
  {
    label: 'Logs',
    icon: 'i-lucide-git-fork',
    to: '/logs',
    active: route.path.startsWith('/logs'),
  },
])

const modal = ref(false)

const itemsDropdown = ref<DropdownMenuItem[]>([
  [
    {
      label: user.value?.username,
      icon: 'i-lucide-user',
    },
    {
      label: user.value?.admin ? 'admin' : 'user',
      icon: 'i-lucide-shield',
    },
  ],
  [
    {
      label: 'New connection',
      icon: 'i-lucide-radio-tower',
      onClick: () => {
        if(!user.value.admin) return toast.add({ title: 'Only admins can create new connections', icon: 'i-lucide-shield-alert', color: 'error' })
        modal.value = true
      },
    },
  ],
  [{
    label: 'Change Password',
    icon: 'i-lucide-key',
  }, {
    label: 'Logout',
    icon: 'i-lucide-log-out',
    onClick: async () => {
      clearUserSession()
      await navigateTo('/login')
    },
  }],
])

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
  modal.value = false
}
</script>

<template>
  <UHeader>
    <template #title>
      <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="32" />
      Sanchez<span class="text-green-500">DNS</span>
    </template>

    <UNavigationMenu v-if="optionSelected" :items="items" />

    <template #right>
      <USelectMenu v-model="optionSelected" :items="optionsConnection" placeholder="Select a connection" size="sm" class="hidden md:block" />

      <UTooltip text="Open on GitHub">
        <UButton color="neutral" variant="ghost" to="https://github.com/rafinhacuri/SanchezDNS" target="_blank" icon="i-simple-icons-github" aria-label="GitHub" />
      </UTooltip>
      <UDropdownMenu :items="itemsDropdown">
        <UButton icon="i-lucide-user" class="rounded-full" color="neutral" variant="outline" size="sm" />
      </UDropdownMenu>
    </template>

    <template #body>
      <USelectMenu v-model="optionSelected" :items="optionsConnection" placeholder="Select a connection" size="sm" />

      <UNavigationMenu :items="items" orientation="vertical" class="-mx-2.5" />
    </template>

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
  </UHeader>
</template>
