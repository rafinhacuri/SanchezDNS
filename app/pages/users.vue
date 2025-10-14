<script setup lang="ts">
useHead({ title: 'Groups' })

useSeoMeta({ description: 'Groups page' })

defineOgImageComponent('Techs', { title: 'Groups page' })

const { optionSelected } = useConection()

const usersConnection = ref(['rafinhacurig@gmail.com', 'user2@example.com', 'user3@example.com', 'user3@example.com', 'user3@example.com', 'user3@example.com', 'user3@example.com', 'user3@example.com', 'user3@example.com'])
</script>

<template>
  <section class="max-w-7xl mx-auto p-4 w-full">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">
        {{ optionSelected }} Member
      </h1>

      <UButton icon="i-lucide-user-plus" class="text-white" variant="soft" color="primary">
        Add Member
      </UButton>
    </div>

    <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-6 shadow-sm">
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100 mb-4">
        Members
      </h2>
      <ul class="divide-y divide-gray-200 dark:divide-gray-700">
        <li v-for="member in usersConnection" :key="member" class="flex items-center justify-between py-3">
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
  </section>
</template>
