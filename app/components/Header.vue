<script setup lang="ts">
import type { DropdownMenuItem, NavigationMenuItem } from '@nuxt/ui'

const route = useRoute()

const { user, clearUserSession } = useUserSession()

const { isLoading, start, finish } = useLoadingIndicator()

const { optionSelected, optionsConnection, refreshConnections } = await useConnection()

const toast = useToast()

const items = computed<NavigationMenuItem[]>(() => [
  { label: 'Zones', icon: 'i-lucide-database', to: '/zones', active: route.path.startsWith('/zones') },

  ...(user.value?.admin ? [{ label: 'Server', icon: 'i-lucide-server', active: route.path.startsWith('/srv/statistics') || route.path.startsWith('/srv/configuration'), children: [{ label: 'Statistics', icon: 'i-lucide-bar-chart-3', description: 'View server statistics.', to: '/srv/statistics' }, { label: 'Configuration', icon: 'i-lucide-settings', description: 'Manage server connection.', to: '/srv/configuration' }] }] : [{ label: 'Statistics', icon: 'i-lucide-bar-chart-3', active: route.path.startsWith('/srv/statistics'), to: '/srv/statistics' }, { label: 'Configuration', icon: 'i-lucide-settings', active: route.path.startsWith('/srv/configuration'), to: '/srv/configuration' }]),

  ...(user.value?.admin ? [{ label: 'Users', icon: 'i-lucide-users', to: '/users', active: route.path.startsWith('/users') }] : []),

  ...(user.value?.admin ? [{ label: 'Logs', icon: 'i-lucide-git-fork', to: '/logs', active: route.path.startsWith('/logs') }] : []),
])

const itemsAdmin = computed<NavigationMenuItem[]>(() => [
  ...(user.value?.admin ? [{ label: 'Connections', icon: 'i-lucide-wifi', active: route.path.startsWith('/dns-connections'), to: '/dns-connections' }] : []),
])

const modal = ref(false)

const modalPassword = ref(false)

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
        if(route.path.startsWith('/dns-connections')) return toast.add({ title: 'Create new connections on the page', icon: 'i-lucide-shield-alert', color: 'error' })
        modal.value = true
      },
    },
  ],
  [{
    label: 'Change Password',
    icon: 'i-lucide-key',
    onClick: () => {
      modalPassword.value = true
    },
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
  refreshConnections()
  modal.value = false
}

watch(modal, newVal => {
  if(!newVal) state.value = { name: '', host: '', apiKey: '', serverId: '', users: [] }
})

const newPassword = ref('')
const confirmPassword = ref('')
const show = ref(false)

function checkStrength(str: string){
  const requirements = [
    { regex: /.{8,}/, text: 'At least 8 characters' },
    { regex: /\d/, text: 'At least 1 number' },
    { regex: /[a-z]/, text: 'At least 1 lowercase letter' },
    { regex: /[A-Z]/, text: 'At least 1 uppercase letter' },
  ]

  return requirements.map(req => ({ met: req.regex.test(str), text: req.text }))
}

const strength = computed(() => checkStrength(newPassword.value))
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

async function changePassword(){
  start()

  if(newPassword.value.length < 8){
    toast.add({ title: 'Password must be at least 8 characters long', icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  if(newPassword.value !== confirmPassword.value){
    toast.add({ title: 'Passwords do not match', icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/api/user/password', { method: 'PATCH', body: { password: newPassword.value, email: user.value?.username } })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  modalPassword.value = false
}
</script>

<template>
  <UHeader>
    <template #title>
      <NuxtImg src="/logo.png" alt="SanchezDNS Logo" width="32" />
      Sanchez<span class="text-green-500">DNS</span>
    </template>

    <UNavigationMenu v-if="optionSelected" :items="items" />
    <UNavigationMenu v-if="user.admin" :items="itemsAdmin" />

    <template #right>
      <USelectMenu v-model="optionSelected" :items="optionsConnection" label-key="name" value-key="_id" placeholder="Select a connection" size="sm" class="hidden md:block " />

      <UDropdownMenu :items="itemsDropdown">
        <UButton icon="i-lucide-user" class="rounded-full" color="neutral" variant="outline" size="sm" />
      </UDropdownMenu>
    </template>

    <template #body>
      <USelectMenu v-model="optionSelected" :items="optionsConnection" label-key="name" value-key="_id" placeholder="Select a connection" size="sm" />

      <UNavigationMenu :items="items" orientation="vertical" class="-mx-2.5" />
      <UNavigationMenu v-if="user.admin" :items="itemsAdmin" orientation="vertical" class="-mx-2.5" />
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

    <UModal v-model:open="modalPassword" title="Change Password" description="Change your account password" :ui="{ footer: 'justify-end' }">
      <template #body>
        <div class="space-y-2">
          <UFormField label="Password">
            <UInput v-model="newPassword" placeholder="Password" :color="color" :type="show ? 'text' : 'password'" :aria-invalid="score < 4" aria-describedby="password-strength" :ui="{ trailing: 'pe-1' }" class="w-full">
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
        <UFormField label="Confirm Password" class="mt-4">
          <UInput v-model="confirmPassword" class="w-full" icon="i-lucide-lock" :type="show ? 'text' : 'password'">
            <template #trailing>
              <UButton color="neutral" variant="link" size="sm" :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'" :aria-label="show ? 'Hide password' : 'Show password'" :aria-pressed="show" aria-controls="password" @click="show = !show" />
            </template>
          </UInput>
        </UFormField>
      </template>

      <template #footer>
        <UButton label="Cancel" :loading="isLoading" variant="outline" @click="modalPassword = false" />
        <UButton label="Confirm" :loading="isLoading" @click="changePassword" />
      </template>
    </UModal>
  </UHeader>
</template>
